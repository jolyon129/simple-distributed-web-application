package backend

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"zl2501-final-project/backend/pb"
	"zl2501-final-project/backend/constant"
)

func init() {
	// Set global logger
	log.SetPrefix("BE Service: ")
	log.SetFlags(log.Ltime | log.Lshortfile)
}




func StartService() {
	log.Println("Backend Server is going to start at: http://localhost:" + constant.Port)
	lis, err := net.Listen("tcp", ":"+constant.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBackendServer(s,&backendServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
