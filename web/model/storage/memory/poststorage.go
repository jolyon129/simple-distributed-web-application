package memory

import (
	"container/list"
	"sync"
	"time"
	"zl2501-final-project/web/model/storage"
)

type MemPostStore struct {
	sync.Mutex
	postMap   map[uint]*storage.PostEntity // // Map index to entity/record
	posts     *list.List
	pkCounter uint
}

// Return a new primary key.
// No need to lock.
func (m *MemPostStore) getNewPK() uint {
	m.pkCounter++
	return m.pkCounter
}

func (m *MemPostStore) Create(post *storage.PostEntity) (uint, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	pk := m.getNewPK()
	newPost := storage.PostEntity{
		ID:          pk,
		UserID:      post.UserID,
		Content:     post.Content,
		CreatedTime: time.Now(),
	}
	m.postMap[newPost.ID] = &newPost
	m.posts.PushBack(&newPost)
	return pk, nil
}

func (m *MemPostStore) Delete(ID uint) *storage.MyStorageError {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.postMap[ID]; !ok {
		return &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		for e := m.posts.Front(); e != nil; e = e.Next() {
			u := e.Value.(storage.PostEntity)
			if u.ID == ID {
				id := u.ID
				delete(m.postMap, id)
				m.posts.Remove(e)
			}
		}
		return nil
	}
}

func (m *MemPostStore) Read(ID uint) (storage.PostEntity, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	var entity storage.PostEntity
	if _, ok := m.postMap[ID]; !ok {
		return entity, &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		eInDB := m.postMap[ID]
		entity = storage.PostEntity{
			ID:          eInDB.ID,
			UserID:      eInDB.UserID,
			Content:     eInDB.Content,
			CreatedTime: eInDB.CreatedTime,
		}
		return entity, nil
	}
}

//
//func (m *MemPostStore) Update(ID uint, post *storage.PostEntity) (uint, *storage.MyStorageError) {
//	m.Lock()
//	defer m.Unlock()
//	if _, ok := m.postMap[ID]; !ok {
//		return 0, &storage.MyStorageError{Message: "Non-exist ID"}
//	} else {
//		createTime := m.postMap[ID].CreatedTime
//		m.postMap[ID] = &storage.PostEntity{
//			ID:         post.ID,
//			Content:    post.Content,
//			CreatedTime: createTime,
//		}
//	}
//}
