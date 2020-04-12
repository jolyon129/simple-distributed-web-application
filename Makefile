GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# compile proto files and generate gRPC code
proctoc:
	protoc -I commonpb commonpb/backend.proto --go_out=plugins=grpc:backend/pb
	protoc -I commonpb commonpb/auth.proto --go_out=plugins=grpc:auth/pb


vendor-web:
	cd ./web && go mod vendor

run-web: vendor-web
	go run ./cmd/web/web.go

build-web: vendor-web
	mkdir -p ./build
	$(GOBUILD) -o ./build/ ./cmd/web/web.go

test:
	go test -v --race ./...
