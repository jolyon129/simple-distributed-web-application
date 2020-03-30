package storage

import (
	"container/list"
	"time"
)

type UserEntity struct {
	ID        uint // The DB will fill this field
	UserName  string
	Password  string
	Follower  *list.List
	Following *list.List
	Posts     *list.List // The oldest comes first
}

type PostEntity struct {
	ID          uint
	UserID      uint
	Content     string
	CreatedTime time.Time
}

type MyStorageError struct {
	Message string
}

func (e *MyStorageError) Error() string {
	return e.Message
}

func (e *MyStorageError) String() string {
	return e.Message
}

type UserStorageInterface interface {
	Create(user *UserEntity) (uint, *MyStorageError)
	Delete(ID uint) *MyStorageError
	// Read by user ID.
	// Return a copy of user entity.
	Read(ID uint) (*UserEntity, *MyStorageError)
	Update(ID uint, user *UserEntity) (uint, *MyStorageError)
	FindAll() []*UserEntity
}

type PostStorageInterface interface {
	Create(post *PostEntity) (uint, *MyStorageError)
	//Delete(ID uint) *MyStorageError
	// Read by post ID
	// Return a copy of post entity
	Read(ID uint) (PostEntity, *MyStorageError)
	//Update(ID uint, post *PostEntity) (uint,*MyStorageError)

}

var drivers = make(map[string]*Manager)

func RegisterDriver(name string, m *Manager) {
	drivers[name] = m
}

func NewManager(name string) *Manager {
	m := drivers[name]
	return m
}

// A storage manager. The is the entry point for the storage package.
type Manager struct {
	UserStorage UserStorageInterface
	PostStorage PostStorageInterface
}

func (m *Manager) GetUserStorage() UserStorageInterface {
	return m.UserStorage
}

func (m *Manager) GetPostStorage() PostStorageInterface {
	return m.PostStorage
}
