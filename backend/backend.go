package backend

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"zl2501-final-project/backend/constant"
	"zl2501-final-project/backend/pb"
)

func init() {
	// Set global logger
	log.SetPrefix("Backend Service: ")
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func StartService() {
	log.Println("Backend Server is going to start at: http://localhost:" + constant.Port)
	lis, err := net.Listen("tcp", ":"+constant.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	loggerHandler := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		log.Printf("received RPC request: %s", info.FullMethod)
		if err != nil {
			log.Printf("method %q failed: %s", info.FullMethod, err)
		}
		return resp, err
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(loggerHandler))
	pb.RegisterBackendServer(s,&backendServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
