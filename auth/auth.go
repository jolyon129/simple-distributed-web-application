package auth

import (
    "google.golang.org/grpc"
    "log"
    "net"
    "zl2501-final-project/auth/constant"
    "zl2501-final-project/auth/pb"
)

func init() {
    // Set global logger
    log.SetPrefix("BE Service: ")
    log.SetFlags(log.Ltime | log.Lshortfile)
}

type server struct {
    pb.UnimplementedAuthServer
}

func StartService() {
    log.Println("Auth(Session) Server is going to start at: http://localhost:" + constant.Port)
    lis, err := net.Listen("tcp", ":"+constant.Port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterAuthServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }

}
