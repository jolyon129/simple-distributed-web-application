package memory

import (
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

func (m *MemTweetStore) Create(tweet *storage.TweetEntity, result chan uint, errorChan chan error) {
	m.Lock()
	defer m.Unlock()
	pk := m.getNewPK()
	newPost := storage.TweetEntity{
		ID:          pk,
		UserID:      tweet.UserID,
		Content:     tweet.Content,
		CreatedTime: time.Now(),
	}
	m.tweetMap[newPost.ID] = &newPost
	//	return pk, nil
	result <- pk
}

//func (m *MemTweetStore) Delete(ID uint) *storage.MyStorageError {
//	m.Lock()
//	defer m.Unlock()
//	if _, ok := m.tweetMap[ID]; !ok {
//		return &storage.MyStorageError{Message: "Non-exist ID"}
//	} else {
//		//for e := m.posts.Front(); e != nil; e = e.Next() {
//		//	u := e.Value.(storage.TweetEntity)
//		//	if u.ID == ID {
//		//		id := u.ID
//		//		delete(m.tweetMap, id)
//		//		m.posts.Remove(e)
//		//	}
//		//}
//		delete(m.tweetMap, ID)
//		return nil
//	}
//}

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
