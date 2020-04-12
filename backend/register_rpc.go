package backend

import (
    "context"
    "zl2501-final-project/backend/constant"
    "zl2501-final-project/backend/model"
    "zl2501-final-project/backend/model/repository"
    pb "zl2501-final-project/backend/pb"
)

var tweetRepo *repository.TweetRepo
var userRepo *repository.UserRepo

func init() {
    tweetRepo = model.GetTweetRepo()
    userRepo = model.GetUserRepo()
}

type backendServer struct {
    pb.BackendServer
}


func (b backendServer) NewTweet(ctx context.Context,
        request *pb.NewTweetRequest) (*pb.NewTweetResponse, error) {
    uId := uint(request.UserId)
    tId, err := tweetRepo.SaveTweet(ctx, repository.TweetInfo{
        UserID:  uId,
        Content: request.Content,
    })
    if err != nil {
        return nil, err
    }
    _, err1 := userRepo.AddTweetToUser(ctx, uId, tId)
    if err1 != nil {
        return nil, err1
    }
    return &pb.NewTweetResponse{
        TweetId: uint64(tId),
    }, nil

}

func (b backendServer) TweetSelectById(ctx context.Context,
        request *pb.SelectByIdRequest) (*pb.TweetSelectByIdResponse, error) {
    tId := uint(request.Id)
    tweet, err := tweetRepo.SelectById(ctx, tId)
    if err != nil {
        return nil, err
    }
    return &pb.TweetSelectByIdResponse{
        Msg: &pb.TweetEntity{
            TweetId:     uint64(tweet.ID),
            UserId:      uint64( tweet.UserID ),
            Content:     tweet.Content ,
            CreatedTime: tweet.CreatedTime.Format(constant.TimeFormat),
        },
    }, nil
}

func (b backendServer) NewUser(ctx context.Context,
    request *pb.NewUserRequest) (*pb.NewUserResponse, error) {

    panic("implement me")
}

func (b backendServer) UserSelectByName(ctx context.Context,
    request *pb.UserSelectByNameRequest) (*pb.UserSelectByNameResponse, error) {
    panic("implement me")
}

func (b backendServer) UserSelectById(ctx context.Context,
    request *pb.SelectByIdRequest) (*pb.UserSelectByIdResponse, error) {
    panic("implement me")
}

func (b backendServer) UserAddTweet(ctx context.Context,
    request *pb.UserAddTweetRequest) (*pb.UserAddTweetResponse, error) {
    panic("implement me")
}

func (b backendServer) FindAllUsers(ctx context.Context,
    request *pb.FindAllUsersRequest) (*pb.FindAllUsersResponse, error) {
    panic("implement me")
}

func (b backendServer) UserCheckWhetherFollowing(ctx context.Context, request *pb.UserCheckWhetherFollowingRequest) (*pb.UserCheckWhetherFollowingResponse, error) {
    panic("implement me")
}

func (b backendServer) StartFollowing(ctx context.Context,
    request *pb.StartFollowingRequest) (*pb.StartFollowingResponse, error) {
    panic("implement me")
}

func (b backendServer) StopFollowing(ctx context.Context,
    request *pb.StopFollowingRequest) (*pb.StopFollowingResponse, error) {
    panic("implement me")
}
