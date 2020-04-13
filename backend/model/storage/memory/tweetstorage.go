package memory

import (
    "errors"
    "sync"
    "time"
    "zl2501-final-project/backend/model/storage"
)

type MemTweetStore struct {
    storage.TweetStorageInterface
    sync.Mutex
    tweetMap map[uint]*storage.TweetEntity // // Map index to entity/record
    //posts     *list.List
    pkCounter uint
}

// Return a new primary key.
// No need to lock.
func (m *MemTweetStore) getNewPK() uint {
    m.pkCounter++
    return m.pkCounter
}


func (m *MemTweetStore) DeleteByCreatedTime(timeStamp time.Time, result chan bool, errorChan chan error){
    m.Lock()
    defer m.Unlock()
    for tId,tweet:= range m.tweetMap{
        if tweet.CreatedTime.Equal(timeStamp){
            delete(m.tweetMap, tId)
            result <- true
            return
        }
    }
    errorChan<-errors.New("didn't find the tweet")
    return
}

func (m *MemTweetStore) Create(tweet *storage.TweetEntity, result chan uint,
    errorChan chan error) uint {
    m.Lock()
    defer m.Unlock()
    pk := m.getNewPK()
    newTweet := storage.TweetEntity{
        ID:          pk,
        UserID:      tweet.UserID,
        Content:     tweet.Content,
        CreatedTime: time.Now(),
    }
    m.tweetMap[newTweet.ID] = &newTweet
    //	return pk, nil
    result <- pk
    return pk
}

func (m *MemTweetStore) Delete(ID uint, result chan bool, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if _, ok := m.tweetMap[ID]; !ok {
        errorChan <- &storage.MyStorageError{Message: "Non-exist ID"}
        return
    } else {
        delete(m.tweetMap, ID)
        result <- true
        return
    }
}

func (m *MemTweetStore) Read(ID uint, result chan *storage.TweetEntity, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    var entity storage.TweetEntity
    if _, ok := m.tweetMap[ID]; !ok {
        errorChan <- &storage.MyStorageError{Message: "Non-exist Tweet ID"}
        return
    } else {
        eInDB := m.tweetMap[ID]
        entity = storage.TweetEntity{
            ID:          eInDB.ID,
            UserID:      eInDB.UserID,
            Content:     eInDB.Content,
            CreatedTime: eInDB.CreatedTime,
        }
        //return entity, nil
        result <- &entity
    }
}

//
//func (m *MemTweetStore) Update(ID uint, post *storage.TweetEntity) (uint, *storage.MyStorageError) {
//	m.Lock()
//	defer m.Unlock()
//	if _, ok := m.tweetMap[ID]; !ok {
//		return 0, &storage.MyStorageError{Message: "Non-exist ID"}
//	} else {
//		createTime := m.tweetMap[ID].CreatedTime
//		m.tweetMap[ID] = &storage.TweetEntity{
//			ID:         post.ID,
//			Content:    post.Content,
//			CreatedTime: createTime,
//		}
//	}
//}
