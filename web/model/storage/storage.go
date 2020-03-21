package storage

import (
	"container/list"
	"time"
)

type UserEntity struct {
	ID        uint // The DB will fill this field
	UserName  string
	Password  string
	followers *list.List
	Posts     *list.List
}

type PostEntity struct {
	ID         uint
	Content    string
	CreateTime time.Time
}

type MyStorageError struct {
	Message string
}

func (e *MyStorageError) Error() string {
	return e.Message
}

type UserStorageInterface interface {
	Create(user *UserEntity) (uint, *MyStorageError)
	Delete(ID uint) *MyStorageError
	Read(ID uint) (*UserEntity, *MyStorageError)
	//Update(ID uint, user *UserEntity) (uint,*MyStorageError)
	FindAll() *list.List
}

type PostStorageInterface interface {
	Create(post *PostEntity) (uint, *MyStorageError)
	Delete(ID uint) *MyStorageError
	Read(ID uint) (*PostEntity, *MyStorageError)
	//Update(ID uint, post *PostEntity) (uint,*MyStorageError)
}

var drivers = make(map[string]*Manager)

func RegisterDriver(name string, models *Manager) {
	drivers[name] = models
}

func NewManager(name string) *Manager {
	m := drivers[name]
	return m
}

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
