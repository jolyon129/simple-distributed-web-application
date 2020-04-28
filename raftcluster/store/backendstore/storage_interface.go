package backendstore

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
	Tweets    *list.List // The oldest comes first
}

type TweetEntity struct {
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
	Create(user *UserEntity, result chan uint, errorChan chan error)
	Delete(ID uint, result chan bool, errorChan chan error)
	// Read by user ID. Return a copy of user entity.
	Read(ID uint, result chan *UserEntity, errorChan chan error)
	Update(ID uint, user *UserEntity, result chan uint, errorChan chan error)
	FindAll(result chan []*UserEntity, errorChan chan error)
	AddTweetToUserDB(uId uint, pId uint, result chan bool, errorChan chan error)
	CheckWhetherFollowingDB(srcId uint, targetId uint, result chan bool, errChan chan error)
	StartFollowingDB(srcId uint, targetID uint, result chan bool, errorChan chan error)
	StopFollowingDB(srcId uint, targetID uint, result chan bool, errorChan chan error)
}

type TweetStorageInterface interface {
	// Return tweet ID
	Create(tweet *TweetEntity, result chan uint, errorChan chan error) uint
	// Return a copy of post entity
	Read(ID uint, result chan *TweetEntity, errorChan chan error)
	Delete(ID uint, result chan bool, errorChan chan error)
	DeleteByCreatedTime(timeStamp time.Time, result chan bool, errorChan chan error)
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
	UserStorage  UserStorageInterface
	TweetStorage TweetStorageInterface
}

//func (m *Manager) GetUserStorage() UserStorageInterface {
//	return m.UserStorage
//}
//
//func (m *Manager) GetTweetStorage() TweetStorageInterface {
//	return m.TweetStorage
//}
