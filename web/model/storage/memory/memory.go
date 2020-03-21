package memory

import (
	"container/list"
	"sync"
	"time"
	"zl2501-final-project/web/model/storage"
)

func init() {
	var memPostStore storage.PostStorageInterface
	var memUserStore storage.UserStorageInterface
	memPostStore = &MemPostStore{postMap: make(map[uint]*storage.PostEntity), pkCounter: 100}
	memUserStore = &MemUserStore{
		userMap:     make(map[uint]*storage.UserEntity),
		users:       list.New(),
		userNameSet: make(map[string]bool),
		pkCounter:   100,
	}
	memModels := storage.Manager{
		UserStorage: memUserStore,
		PostStorage: memPostStore,
	}
	// Register the implementation of storage manager.
	storage.RegisterDriver("memory", &memModels)
}

type MemUserStore struct {
	sync.Mutex
	userMap     map[uint]*storage.UserEntity
	users       *list.List
	userNameSet map[string]bool
	pkCounter   uint
}

func (m *MemUserStore) FindAll() *list.List {
	return m.users
}

type MemPostStore struct {
	sync.Mutex
	postMap   map[uint]*storage.PostEntity
	pkCounter uint
}

func (m *MemUserStore) getNewPK() uint {
	m.pkCounter++
	return m.pkCounter
}

func (m *MemPostStore) getNewPK() uint {
	m.Lock()
	defer m.Unlock()
	m.pkCounter++
	return m.pkCounter
}

func (m *MemUserStore) Create(user *storage.UserEntity) (uint, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.userNameSet[user.UserName]; ok {
		return 0, &storage.MyStorageError{Message: "Duplicate UserStorage Names!"}
	}
	pk := m.getNewPK()
	newUser := storage.UserEntity{
		ID:       pk,
		UserName: user.UserName,
		Password: user.Password,
		Posts:    list.New(),
	}
	m.userMap[pk] = &newUser
	m.userNameSet[user.UserName] = true
	m.users.PushBack(&newUser)
	return pk, nil
}

func (m *MemUserStore) Delete(ID uint) *storage.MyStorageError {
	if _, ok := m.userMap[ID]; !ok {
		return &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		delete(m.userMap, ID)
		for e := m.users.Front(); e != nil; e = e.Next() {
			u := e.Value.(storage.UserEntity)
			if u.ID == ID {
				name := u.UserName
				id := u.ID
				delete(m.userNameSet, name)
				delete(m.userMap, id)
				m.users.Remove(e)
			}
		}
		return nil
	}
}

func (m *MemUserStore) Read(ID uint) (*storage.UserEntity, *storage.MyStorageError) {
	if _, ok := m.userMap[ID]; !ok {
		return nil, &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		return m.userMap[ID], nil
	}
}

//func (m *MemUserStore) Update(ID uint, user *storage.UserEntity) (uint, *storage.MyStorageError) {
//	if _, ok := m.userMap[ID]; !ok {
//		return 0, &storage.MyStorageError{Message: "Non-exist ID"}
//	} else {
//		m.userMap[ID] = &storage.UserEntity{
//			ID:       ID,
//			UserName: user.UserName,
//			Password: user.Password,
//			Posts:    user.Posts,
//		}
//		return ID, nil
//	}
//}

func (m *MemPostStore) Create(post *storage.PostEntity) (uint, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	pk := m.getNewPK()
	m.postMap[pk] = &storage.PostEntity{
		ID:         pk,
		Content:    post.Content,
		CreateTime: time.Now(),
	}
	return pk, nil
}

func (m *MemPostStore) Delete(ID uint) *storage.MyStorageError {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.postMap[ID]; !ok {
		return &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		delete(m.postMap, ID)
		return nil
	}
}

func (m *MemPostStore) Read(ID uint) (*storage.PostEntity, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.postMap[ID]; !ok {
		return nil, &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		return m.postMap[ID], nil
	}
}

//
//func (m *MemPostStore) Update(ID uint, post *storage.PostEntity) (uint, *storage.MyStorageError) {
//	m.Lock()
//	defer m.Unlock()
//	if _, ok := m.postMap[ID]; !ok {
//		return 0, &storage.MyStorageError{Message: "Non-exist ID"}
//	} else {
//		createTime := m.postMap[ID].CreateTime
//		m.postMap[ID] = &storage.PostEntity{
//			ID:         post.ID,
//			Content:    post.Content,
//			CreateTime: createTime,
//		}
//	}
//}
