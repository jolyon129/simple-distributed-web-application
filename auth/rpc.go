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
    pb.UnimplementedAuthServer
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

// Read session from given sessId If its legal.
// If not exist, create a new newSessionId and return.
// If exist and the newSessionId is valid, reuse the same session and return the same one.
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

func (b authService) SetValue(ctx context.Context, req *pb.SetValueRequest) (
        *pb.SetValueResponse, error) {
    ok, err := globalSessions.SetValue(ctx, req.Ssid, req.Key, req.Value)
    if err != nil {
        return nil, err
    }
    return &pb.SetValueResponse{
        Ok: ok,
    }, nil
}

func (b authService) GetValue(ctx context.Context, req *pb.GetValueRequest) (
        *pb.GetValueResponse, error) {
    value, err := globalSessions.GetValue(ctx, req.Ssid, req.Key)
    if err != nil {
        return nil, err
    }
    retstr, _ := value.(string)
    return &pb.GetValueResponse{
        Value: retstr,
    }, nil
}

func (b authService) DeleteValue(ctx context.Context, req *pb.DeleteValueRequest) (*pb.DeleteValueResponse, error) {
    ok, err := globalSessions.DeleteValue(ctx, req.Ssid, req.Key)
    if err != nil {
        return nil, err
    }
    return &pb.DeleteValueResponse{
        Ok: ok,
    }, nil
}
