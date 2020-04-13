package pb

import (
    "google.golang.org/grpc"
    "log"
    "zl2501-final-project/web/constant"
)

var BackendClientIns BackendClient
var AuthClientIns AuthClient

func CreateBackendServiceConnection() *grpc.ClientConn {
    beConn, err := grpc.Dial(constant.BackendServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    // close connection when close the web service
    BackendClientIns = NewBackendClient(beConn)
    return beConn
}

func CreateAuthServiceConnection() *grpc.ClientConn {
    con, err := grpc.Dial(constant.AuthServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    // close connection when close the web service
    BackendClientIns = NewBackendClient(con)
    AuthClientIns = NewAuthClient(con)
    return con
}
