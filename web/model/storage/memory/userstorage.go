package memory

import (
	"container/list"
	"sort"
	"sync"
	"zl2501-final-project/web/model/storage"
)

type MemUserStore struct {
	sync.Mutex
	userMap map[uint]*storage.UserEntity // Map index to entity/record
	//users       *list.List                   // A list of user entity. The entry is a pointer.
	userNameSet map[string]bool // A set of username. Used as a fast approach to avoid duplicate names
	pkCounter   uint            // Primary Key Counter
}

// Sort User By their ID
type EntityById []*storage.UserEntity

func (e EntityById) Len() int {
	return len(e)
}

func (e EntityById) Less(i, j int) bool {
	return e[i].ID < e[j].ID
}

func (e EntityById) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (m *MemUserStore) FindAll() []*storage.UserEntity {
	m.Lock()
	defer m.Unlock()
	//newList := list.New()
	res := make(EntityById, 0)
	for _, entity := range m.userMap {
		newE := storage.UserEntity{}
		copyUserEntity(&newE, entity)
		//newList.PushBack(&newE)
		res = append(res, &newE)
	}
	sort.Sort(res)
	return res
}

// Return a new primary key
// This function does not need to be locked
func (m *MemUserStore) getNewPK() uint {
	m.pkCounter++
	return m.pkCounter
}

// Update can only modified the password and the post, the following and the follower list.
// Take O(#post+#following+#follower) time.
func (m *MemUserStore) Update(ID uint, newUserE *storage.UserEntity) (uint, *storage.MyStorageError) {
	m.Lock()
	defer m.Unlock()
	totalNewE := &storage.UserEntity{}
	copyUserEntity(totalNewE, newUserE)
	m.userMap[ID] = totalNewE
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
		ID:        pk,
		UserName:  user.UserName,
		Password:  user.Password,
		Posts:     list.New(),
		Following: list.New(),
		Follower:  list.New(),
	}
	m.userMap[pk] = &newUser
	m.userNameSet[user.UserName] = true
	//m.users.PushBack(&newUser)
	return pk, nil
}

func (m *MemUserStore) Delete(ID uint) *storage.MyStorageError {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.userMap[ID]; !ok {
		return &storage.MyStorageError{Message: "Non-exist ID"}
	} else {
		userE, _ := m.Read(ID)
		delete(m.userNameSet, userE.UserName)
		delete(m.userMap, ID)
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

// Copy a list of uint
func copyUintList(dst *list.List, src *list.List) {
	if dst == nil {
		dst = list.New()
	}
	if src == nil {
		src = list.New()
	}
	for e := src.Front(); e != nil; e = e.Next() {
		pId := e.Value.(uint)
		dst.PushBack(pId)
	}
}

func copyUserEntity(dst *storage.UserEntity, src *storage.UserEntity) {
	dst.Posts = list.New()
	dst.Following = list.New()
	dst.Follower = list.New()
	copyUintList(dst.Posts, src.Posts)
	copyUintList(dst.Follower, src.Follower)
	copyUintList(dst.Following, src.Following)
	dst.Password = src.Password
	dst.UserName = src.UserName
	dst.ID = src.ID
}
