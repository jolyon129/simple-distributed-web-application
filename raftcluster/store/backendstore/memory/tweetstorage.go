package memory

import (
    "encoding/json"
    "errors"
    "sync"
    "time"
    "zl2501-final-project/raftcluster/store/backendstore"
)

type MemTweetStore struct {
    backendstore.TweetStorageInterface
    sync.Mutex
    tweetMap map[uint]*backendstore.TweetEntity // // Map index to entity/record
    //posts     *list.List
    pkCounter uint
}

// Return a new primary key.
// No need to lock.
func (m *MemTweetStore) getNewPK() uint {
    m.pkCounter++
    return m.pkCounter
}

func (m *MemTweetStore) DeleteByCreatedTime(timeStamp time.Time, result chan bool, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    for tId, tweet := range m.tweetMap {
        if tweet.CreatedTime.Equal(timeStamp) {
            delete(m.tweetMap, tId)
            result <- true
            return
        }
    }
    errorChan <- errors.New("didn't find the tweet")
    return
}

func (m *MemTweetStore) Create(tweet *backendstore.TweetEntity, result chan uint,
        errorChan chan error) uint {
    m.Lock()
    defer m.Unlock()
    pk := m.getNewPK()
    newTweet := backendstore.TweetEntity{
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
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist ID"}
        return
    } else {
        delete(m.tweetMap, ID)
        result <- true
        return
    }
}

func (m *MemTweetStore) Read(ID uint, result chan *backendstore.TweetEntity, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    var entity backendstore.TweetEntity
    if _, ok := m.tweetMap[ID]; !ok {
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist Tweet ID"}
        return
    } else {
        eInDB := m.tweetMap[ID]
        entity = backendstore.TweetEntity{
            ID:          eInDB.ID,
            UserID:      eInDB.UserID,
            Content:     eInDB.Content,
            CreatedTime: eInDB.CreatedTime,
        }
        //return entity, nil
        result <- &entity
    }
}

// Get the snapshot of the tweetMap
func (m *MemTweetStore) GetSnapshot() ([]byte, error) {
    return m.MarshalJSON()
}

func (m *MemTweetStore) MarshalJSON() ([]byte, error) {
    m.Lock()
    defer m.Unlock()
    return json.Marshal(map[string]interface{}{
        "pkCounter": m.pkCounter,
        "tweetMap":  m.tweetMap,
    })
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
