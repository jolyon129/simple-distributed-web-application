package store

// This file is copied and edited from the backend repo module.
// I don't want to rewrite a the same logic in order to call the functions inside backend store
// The wrappers help to provide context and build result and error channel for storage functions

import (
    "context"
    "time"
    beStorage "zl2501-final-project/raftcluster/store/backendstore"
    "zl2501-final-project/web/constant"
)

// Create a new user and return user id
// if the user name is not duplicated.
// Otherwise return error
func CreateNewUser(ctx context.Context, u *UserInfo) (uint, error) {
    result := make(chan uint)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.UserStorage.Create(&beStorage.UserEntity{
            ID:       0,
            UserName: u.UserName,
            Password: u.Password,
        }, result, errorChan)
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return 0, err
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}

//func UserSelectByName(ctx context.Context, name string) (*beStorage.UserEntity,
//        error) {
//    result := make(chan []*beStorage.UserEntity)
//    errorChan := make(chan error)
//    go func() {
//        bkStorageManager.UserStorage.FindAll(result, errorChan)
//    }()
//    select {
//    case users := <-result:
//        for _, value := range users {
//            if value.UserName == name {
//                return value, nil
//            }
//        }
//        return nil, errors.New("the name does not exist")
//    case err := <-errorChan:
//        return nil, err
//    case <-ctx.Done():
//        return nil, ctx.Err()
//    }
//}

func UserSelectById(ctx context.Context, uid uint) (*beStorage.UserEntity, error) {
    result := make(chan *beStorage.UserEntity)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.UserStorage.Read(uid, result, errorChan)
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return nil, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// Add the tweet id into the user
func AddTweetToUser(ctx context.Context, uId uint, pId uint) (bool, error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.UserStorage.AddTweetToUserDB(uId, pId, result, errorChan)
    }()
    select {
    case <-result:
        return true, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// Return all users in the database
func FindAllUsers(ctx context.Context) ([]*beStorage.UserEntity, error) {
    result := make(chan []*beStorage.UserEntity)
    errorChan := make(chan error)
    go func() {
        bkStorageManager.UserStorage.FindAll(result, errorChan)
    }()
    select {
    case users := <-result:
        return users, nil
    case err := <-errorChan:
        return nil, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// Check whether the user srcId follows the user targetId.
// Take O(#following) time
func CheckWhetherFollowing(ctx context.Context, srcId uint, targetId uint) (bool,
        error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.UserStorage.CheckWhetherFollowingDB(srcId, targetId, result, errorChan)
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// User srcId starts to follow targetId.
func StartFollowing(ctx context.Context, srcId uint, targetId uint) (bool, error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.UserStorage.StartFollowingDB(srcId, targetId, result, errorChan)
    }()
    select {
    case err := <-errorChan:
        return false, err
    case ret := <-result:
        return ret, nil
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// srcId stop following targetId.
// targetId remove the follower srcId.
func StopFollowing(ctx context.Context, srcId uint, targetId uint) (bool, error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.UserStorage.StopFollowingDB(srcId, targetId, result, errorChan)
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

func SaveTweet(ctx context.Context, p TweetInfo) (uint, error) {
    result := make(chan uint)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    newTweet := &beStorage.TweetEntity{
        UserID:      p.UserID,
        Content:     p.Content,
        CreatedTime: time.Time{}}
    go func() {
        bkStorageManager.TweetStorage.Create(newTweet, result, errorChan)
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return 0, err
    case <-ctx.Done():
        ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        TweetDeleteByCreatedTime(ctx1,
            newTweet.CreatedTime) // try to delete teh already created tweet.
        return 0, ctx.Err()
    }
}

// Return a post Entity according the post id
func TweetSelectById(ctx context.Context, pId uint) (*beStorage.TweetEntity, error) {
    result := make(chan *beStorage.TweetEntity)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.TweetStorage.Read(pId, result, errorChan)
    }()
    select {
    case tweet := <-result:
        return tweet, nil
    case err := <-errorChan:
        return nil, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// Delete if the tweet existed
func DeleteById(ctx context.Context, tId uint) (bool, error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.TweetStorage.Delete(tId, result, errorChan)
    }()
    select {
    case ok := <-result:
        return ok, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

func TweetDeleteByCreatedTime(ctx context.Context, timestamp time.Time) (bool,
        error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        bkStorageManager.TweetStorage.DeleteByCreatedTime(timestamp, result, errorChan)
    }()
    select {
    case ok := <-result:
        return ok, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}
