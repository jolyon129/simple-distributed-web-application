package backend

import (
	"container/list"
	"context"
	"log"
	"zl2501-final-project/backend/constant"
	"zl2501-final-project/backend/model"
	"zl2501-final-project/backend/model/repository"
	backendpb "zl2501-final-project/backend/pb"
)

var tweetRepo *repository.TweetRepo
var userRepo *repository.UserRepo

func init() {
	tweetRepo = model.GetTweetRepo()
	userRepo = model.GetUserRepo()
}

type backendServer struct {
	backendpb.UnimplementedBackendServer
}

func (b backendServer) NewTweet(ctx context.Context,
	request *backendpb.NewTweetRequest) (*backendpb.NewTweetResponse, error) {
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
	return &backendpb.NewTweetResponse{
		TweetId: uint64(tId),
	}, nil

}

func (b backendServer) TweetSelectById(ctx context.Context,
	request *backendpb.SelectByIdRequest) (*backendpb.TweetSelectByIdResponse, error) {
	tId := uint(request.Id)
	tweet, err := tweetRepo.SelectById(ctx, tId)
	if err != nil {
		return nil, err
	}
	return &backendpb.TweetSelectByIdResponse{
		Msg: &backendpb.TweetEntity{
			TweetId:     uint64(tweet.ID),
			UserId:      uint64(tweet.UserID),
			Content:     tweet.Content,
			CreatedTime: tweet.CreatedTime.Format(constant.TimeFormat),
		},
	}, nil
}

func (b backendServer) NewUser(ctx context.Context,
	request *backendpb.NewUserRequest) (*backendpb.NewUserResponse, error) {
	uId, err := userRepo.CreateNewUser(ctx, &repository.UserInfo{
		UserName: request.UserName,
		Password: request.UserPwd,
	})
	if err != nil {
		return nil, err
	}
	return &backendpb.NewUserResponse{
		UserId: uint64(uId),
	}, nil
}

func (b backendServer) UserSelectByName(ctx context.Context,
	request *backendpb.UserSelectByNameRequest) (*backendpb.UserSelectByNameResponse, error) {
	user, err := userRepo.SelectByName(ctx, request.Name)
	//test,err1:=userRepo.SelectById(ctx,1)
	if err != nil {
		return nil, err
	}
	return &backendpb.UserSelectByNameResponse{
		User: &backendpb.UserEntity{
			UserId:     uint64(user.ID),
			UserName:   user.UserName,
			Password:   "",
			Followers:  convertUintListToUint64Slice(user.Follower),
			Followings: convertUintListToUint64Slice(user.Following),
			Tweets:     convertUintListToUint64Slice(user.Tweets),
		},
	}, nil
}

func convertUintListToUint64Slice(list *list.List) []uint64 {
	tweets := make([]uint64, list.Len())
	i := 0
	for e := list.Back(); e != nil; e = e.Prev() {
		tid := e.Value.(uint)
		tweets[i] = uint64(tid)
		i++
	}
	return tweets
}

func (b backendServer) UserSelectById(ctx context.Context,
	request *backendpb.SelectByIdRequest) (*backendpb.UserSelectByIdResponse, error) {
	user, err := userRepo.SelectById(ctx, uint(request.Id))
	if err != nil {
		return nil, err
	}
	return &backendpb.UserSelectByIdResponse{
		User: &backendpb.UserEntity{
			UserId:     uint64(user.ID),
			UserName:   user.UserName,
			Password:   "",
			Followers:  convertUintListToUint64Slice(user.Follower),
			Followings: convertUintListToUint64Slice(user.Following),
			Tweets:     convertUintListToUint64Slice(user.Tweets),
		},
	}, nil
}

func (b backendServer) UserAddTweet(ctx context.Context,
	request *backendpb.UserAddTweetRequest) (*backendpb.UserAddTweetResponse, error) {
	uId := uint(request.UserId)
	_, err := userRepo.SelectById(ctx, uId)
	if err != nil {
		return nil, err
	}
	tId, err1 := tweetRepo.SaveTweet(ctx, repository.TweetInfo{
		UserID:  uId,
		Content: request.Content,
	})
	if err1 != nil {
		return nil, err
	}
	_, err2 := userRepo.AddTweetToUser(ctx, uId, tId)
	if err2 != nil {
		return nil, err2
	}
	return &backendpb.UserAddTweetResponse{
		Ok: true,
	}, nil
}

func (b backendServer) FindAllUsers(ctx context.Context,
	request *backendpb.FindAllUsersRequest) (*backendpb.FindAllUsersResponse, error) {
	users, err := userRepo.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	retUsers := make([]*backendpb.UserEntity, len(users))
	for idx, user := range users {
		retUsers[idx] = &backendpb.UserEntity{
			UserId:     uint64(user.ID),
			UserName:   user.UserName,
			Password:   "",
			Followers:  convertUintListToUint64Slice(user.Follower),
			Followings: convertUintListToUint64Slice(user.Following),
			Tweets:     convertUintListToUint64Slice(user.Tweets),
		}
	}
	return &backendpb.FindAllUsersResponse{
		Users: retUsers,
	}, nil
}

func (b backendServer) UserCheckWhetherFollowing(ctx context.Context,
	request *backendpb.UserCheckWhetherFollowingRequest) (*backendpb.UserCheckWhetherFollowingResponse, error) {
	srcId := uint(request.SourceUserId)
	tarId := uint(request.TargetUserId)
	ok, err := userRepo.CheckWhetherFollowing(ctx, srcId, tarId)
	if err != nil {
		return nil, err
	}
	return &backendpb.UserCheckWhetherFollowingResponse{
		Ok: ok,
	}, nil
}

func (b backendServer) StartFollowing(ctx context.Context,
	request *backendpb.StartFollowingRequest) (*backendpb.StartFollowingResponse, error) {
	srcId := uint(request.SourceUserId)
	tarId := uint(request.TargetUserId)
	ok, err := userRepo.StartFollowing(ctx, srcId, tarId)
	if err != nil {
		return nil, err
	}
	return &backendpb.StartFollowingResponse{
		Ok: ok,
	}, nil
}

func (b backendServer) StopFollowing(ctx context.Context,
	request *backendpb.StopFollowingRequest) (*backendpb.StopFollowingResponse, error) {
	srcId := uint(request.SourceUserId)
	tarId := uint(request.TargetUserId)
	ok, err := userRepo.StopFollowing(ctx, srcId, tarId)
	if err != nil {
		return nil, err
	}
	return &backendpb.StopFollowingResponse{
		Ok: ok,
	}, nil
}
func (s *backendServer) SayHello(ctx context.Context, in *backendpb.HelloRequest) (*backendpb.HelloReply,
	error) {
	log.Printf("Received: %v", in.GetName())
	return &backendpb.HelloReply{Message: "Hello " + in.GetName()}, nil
}