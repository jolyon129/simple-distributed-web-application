package repository

import (
    "context"
    "sync"
    "time"
    "zl2501-final-project/backend/constant"
    "zl2501-final-project/backend/model/storage"
)

type TweetRepo struct {
    sync.Mutex
    storage storage.TweetStorageInterface
}

type TweetInfo struct {
    UserID  uint
    Content string
}

func (postRepo *TweetRepo) SaveTweet(ctx context.Context, p TweetInfo) (uint, error) {
    result := make(chan uint, 1)
    errorChan := make(chan error, 1)
    newTweet := &storage.TweetEntity{
        UserID:      p.UserID,
        Content:     p.Content,
        CreatedTime: time.Time{}}
    go func() {
        postRepo.storage.Create(newTweet, result, errorChan)
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return 0, err
    case <-ctx.Done():
        ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        postRepo.DeleteByCreatedTime(ctx1,
            newTweet.CreatedTime) // try to delete teh already created tweet.
        return 0, ctx.Err()
    }
}

// Return a post Entity according the post id
func (postRepo *TweetRepo) SelectById(ctx context.Context, pId uint) (*storage.TweetEntity, error) {
    result := make(chan *storage.TweetEntity, 1)
    errorChan := make(chan error, 1)
    go func() {
        postRepo.storage.Read(pId, result, errorChan)
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
func (postRepo *TweetRepo) DeleteById(ctx context.Context, tId uint) (bool, error) {
    result := make(chan bool, 1)
    errorChan := make(chan error, 1)
    go func() {
        postRepo.storage.Delete(tId, result, errorChan)
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

func (postRepo *TweetRepo) DeleteByCreatedTime(ctx context.Context, timestamp time.Time) (bool,
        error) {
    result := make(chan bool, 1)
    errorChan := make(chan error, 1)
    go func() {
        postRepo.storage.DeleteByCreatedTime(timestamp, result, errorChan)
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

func NewTweetRepo(storageInterface storage.TweetStorageInterface) *TweetRepo {
    ret := TweetRepo{storage: storageInterface}
    return &ret
}
