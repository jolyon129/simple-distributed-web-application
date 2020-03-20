package memory

import (
	"sync"
	"zl2501-final-project/web/models"
	"zl2501-final-project/web/models/storage"
)

func init() {
	var memPostStore storage.PostStoreInterface
	var memUserStore storage.UserStoreInterface
	memPostStore = &MemPostStore{Posts: make(map[int]storage.Post)}
	memUserStore = &MemUserStore{Users: make(map[int]storage.User)}
	memModels := models.Models{
		UserStore: &memUserStore,
		PostStore: &memPostStore,
	}
	models.RegisterDriver("memory", memModels)
}

type MemUserStore struct {
	sync.Mutex
	Users     map[int]storage.User
	pkCounter uint
}

func (m *MemUserStore) Create(user *storage.User) int {
	m.Lock()
	defer m.Unlock()
}

func (m *MemUserStore) Delete(user *storage.User) {
	panic("implement me")
}

func (m *MemUserStore) Get(ID int) *storage.User {
	panic("implement me")
}

func (m *MemUserStore) Update(ID int, user *storage.User) int {
	panic("implement me")
}

type MemPostStore struct {
	Posts map[int]storage.Post
}

func (m *MemPostStore) Create(post *storage.Post) {
	panic("implement me")
}

func (m *MemPostStore) Delete(post *storage.Post) {
	panic("implement me")
}

func (m *MemPostStore) Get(ID int) *storage.Post {
	panic("implement me")
}

func (m *MemPostStore) Update(ID int, post *storage.Post) int {
	panic("implement me")
}
