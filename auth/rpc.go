package auth

import (
    "context"
    "zl2501-final-project/auth/pb"
    "zl2501-final-project/auth/sessmanager"
)
var globalSessions *sessmanager.Manager

func init() {
    globalSessions, _ = sessmanager.GetManagerSingleton("memory")
}

type authService struct {
    pb.AuthServer
}

func (b authService) SessionAuth(ctx context.Context,
        request *pb.SessionGeneralRequest) (*pb.SessionAuthResponse, error) {
    ok, err := globalSessions.SessionAuth(ctx, request.SessionId)
    if err != nil {
        return nil, err
    }
    return &pb.SessionAuthResponse{
        Ok: ok,
    }, nil

}

func (b authService) SessionStart(ctx context.Context,
        request *pb.SessionGeneralRequest) (*pb.SessionGeneralResponse, error) {
    nSessId, err := globalSessions.SessionStart(ctx, request.SessionId)
    if err != nil {
        return nil, err
    }
    return &pb.SessionGeneralResponse{
        SessionId: nSessId,
    }, nil
}

func (b authService) SessionDestroy(ctx context.Context,
        request *pb.SessionGeneralRequest) (*pb.SessionDestroyResponse, error) {
    ok, err := globalSessions.SessionDestroy(ctx, request.SessionId)
    if err != nil {
        return nil, err
    }
    return &pb.SessionDestroyResponse{
        Ok: ok,
    }, nil
}

