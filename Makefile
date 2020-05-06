GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# compile proto files and generate gRPC code
proctoc:
	protoc -I commonpb commonpb/backend.proto --go_out=plugins=grpc:web/pb
	protoc -I commonpb commonpb/backend.proto --go_out=plugins=grpc:backend/pb
	protoc -I commonpb commonpb/auth.proto --go_out=plugins=grpc:web/pb
	protoc -I commonpb commonpb/auth.proto --go_out=plugins=grpc:auth/pb


vendor-web:
	cd ./web && go mod vendor

vendor-backend:
	cd ./backend && go mod vendor

vendor-auth:
	cd ./auth && go mod vendor

vendor-raftcluster:
	cd ./raftcluster && go mod vendor

vendor-all: vendor-backend vendor-auth vendor-web vendor-raftcluster

run-web: vendor-web
	go run ./cmd/web/web.go

build-web: vendor-web
	mkdir -p ./build
	$(GOBUILD) -o ./build/ ./cmd/web/web.go

build-auth: vendor-auth
	mkdir -p ./build
	$(GOBUILD) -o ./build/ ./cmd/auth/auth.go

build-backend: vendor-backend
	mkdir -p ./build
	$(GOBUILD) -o ./build/ ./cmd/backend/backend.go

build: vendor-all
	make build-web
	make build-auth
	make build-backend

run-backend: vendor-backend
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run ./cmd/backend/backend.go


run-auth: vendor-auth
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run ./cmd/auth/auth.go


run-raft: vendor-raftcluster
	go run ./cmd/raftcluster/raftcluster.go

test: vendor-all
	go run ./cmd/raftcluster/raftcluster.go
	go test -v --race ./auth/...
	go test -v --race ./backend/...
	#go test -v --race ./raftcluster/store/...


