package memory

import (
	"container/list"
	"sync"
	"zl2501-final-project/web/model/storage"
)

//TODO:
// return a new user entity to avoid clients change the
// internal data

type MemUserStore struct {
	sync.Mutex
	userMap     map[uint]*storage.UserEntity // Map index to entity/record
	users       *list.List                   // A list of user entity. The entry is a pointer.
	userNameSet map[string]bool              // A set of username. Used as a fast approach to avoid duplicate names
	pkCounter   uint                         // Primary Key Counter
}

func (m *MemUserStore) FindAll() *list.List {
	m.Lock()
	defer m.Unlock()
	newList := list.New()
	for e := m.users.Front(); e != nil; e = e.Next() {
		u := e.Value.(*storage.UserEntity)
		newUE := storage.UserEntity{}
		copyUserEntity(&newUE, u)
		newList.PushBack(&newUE)
	}
	return newList
}

// Return a new primary key
// This function does not need to be locked
func (m *MemUserStore) getNewPK() uint {
	m.pkCounter++
	return m.pkCounter
}

// Update can only modified the password and the post list.
// Take O(#post) time.
func (m *MemUserStore) Update(ID uint, user *storage.UserEntity) (uint, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	userEntity := m.userMap[ID]
	userEntity.Password = user.Password
	newPostList := list.New()
	// Copy the post list
	copyPostList(newPostList, user.Posts)
	userEntity.Posts = newPostList // Change the pointer in the userEntity
	return ID, nil
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
	m.Lock()
	defer m.Unlock()
	if _, ok := m.userMap[ID]; !ok {
		return &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
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
	m.Lock()
	defer m.Unlock()

	if _, ok := m.userMap[ID]; !ok {
		return nil, &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		uInDB := m.userMap[ID]
		// Copy the post list
		newUser := storage.UserEntity{}
		copyUserEntity(&newUser, uInDB)
		return &newUser, nil
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

func copyPostList(dst *list.List, src *list.List) {
	for e := src.Front(); e != nil; e = e.Next() {
		pId := e.Value.(uint)
		dst.PushBack(pId)
	}
}

func copyUserEntity(dst *storage.UserEntity, src *storage.UserEntity) {
	dst.Posts = list.New()
	copyPostList(dst.Posts, src.Posts)
	dst.Password = src.Password
	dst.UserName = src.UserName
	dst.ID = src.ID
}
