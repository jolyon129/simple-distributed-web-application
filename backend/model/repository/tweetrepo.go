package repository

import (
    "context"
    "sync"
    "time"
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
    go func() {
        postRepo.storage.Create(&storage.TweetEntity{
            ID:          0,
            UserID:      p.UserID,
            Content:     p.Content,
            CreatedTime: time.Time{},
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

func NewTweetRepo(storageInterface storage.TweetStorageInterface) *TweetRepo {
    ret := TweetRepo{storage: storageInterface}
    return &ret
}
