package backend

import (
    "context"
    "zl2501-final-project/backend/session/sessmanager"
    pb "zl2501-final-project/internal/backendpb"
)
var globalSessions *sessmanager.Manager

func init() {
    globalSessions, _ = sessmanager.GetManagerSingleton("memory")
}

type backendServer struct {
    pb.BackendServer
}

func (b backendServer) SessionAuth(ctx context.Context,
        request *pb.SessionGeneralRequest) (*pb.SessionAuthResponse, error) {
    ok, err := globalSessions.SessionAuth(ctx, request.SessionId)
    if err != nil {
        return nil, err
    }
    return &pb.SessionAuthResponse{
        Ok: ok,
    }, nil

}

func (b backendServer) SessionStart(ctx context.Context,
        request *pb.SessionGeneralRequest) (*pb.SessionGeneralResponse, error) {
    nSessId, err := globalSessions.SessionStart(ctx, request.SessionId)
    if err != nil {
        return nil, err
    }
    return &pb.SessionGeneralResponse{
        SessionId: nSessId,
    }, nil
}

func (b backendServer) SessionDestroy(ctx context.Context, request *pb.SessionGeneralRequest) (*pb.SessionDestroyResponse, error) {
    ok, err := globalSessions.SessionDestroy(ctx, request.SessionId)
    if err != nil {
        return nil, err
    }
    return &pb.SessionDestroyResponse{
        Ok: ok,
    }, nil
}

func (b backendServer) NewTweet(ctx context.Context, request *pb.NewTweetRequest) (*pb.NewTweetResponse, error) {
    panic("implement me")
}

func (b backendServer) TweetSelectById(ctx context.Context, request *pb.SelectByIdRequest) (*pb.TweetSelectByIdResponse, error) {
    panic("implement me")
}

func (b backendServer) NewUser(ctx context.Context, request *pb.NewUserRequest) (*pb.NewUserResponse, error) {

    panic("implement me")
}

func (b backendServer) UserSelectByName(ctx context.Context, request *pb.UserSelectByNameRequest) (*pb.UserSelectByNameResponse, error) {
    panic("implement me")
}

func (b backendServer) UserSelectById(ctx context.Context, request *pb.SelectByIdRequest) (*pb.UserSelectByIdResponse, error) {
    panic("implement me")
}

func (b backendServer) UserAddTweet(ctx context.Context, request *pb.UserAddTweetRequest) (*pb.UserAddTweetResponse, error) {
    panic("implement me")
}

func (b backendServer) FindAllUsers(ctx context.Context, request *pb.FindAllUsersRequest) (*pb.FindAllUsersResponse, error) {
    panic("implement me")
}

func (b backendServer) UserCheckWhetherFollowing(ctx context.Context, request *pb.UserCheckWhetherFollowingRequest) (*pb.UserCheckWhetherFollowingResponse, error) {
    panic("implement me")
}

func (b backendServer) StartFollowing(ctx context.Context, request *pb.StartFollowingRequest) (*pb.StartFollowingResponse, error) {
    panic("implement me")
}

func (b backendServer) StopFollowing(ctx context.Context, request *pb.StopFollowingRequest) (*pb.StopFollowingResponse, error) {
    panic("implement me")
}
