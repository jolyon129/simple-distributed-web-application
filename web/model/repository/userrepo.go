package repository

import (
	"errors"
	"log"
	"sync"
	"zl2501-final-project/web/model/storage"
)

// This is a singleton type
type UserRepo struct {
	sync.Mutex // Need lock
	Storage    storage.UserStorageInterface
	con        *sync.Cond
}

type UserInfo struct {
	UserName string
	Password string
}

func NewUserRepo() *UserRepo {
	ret := UserRepo{
		Storage: nil,
	}
	ret.con = sync.NewCond(&ret)
	return &ret
}

// Create a new user and return user id
// if the user name is no duplicated.
// Otherwise return error
func (userRepo *UserRepo) CreateNewUser(u *UserInfo) (uint, error) {
	userRepo.Lock()
	defer userRepo.Unlock()
	ID, err := userRepo.Storage.Create(&storage.UserEntity{
		ID:       0,
		UserName: u.UserName,
		Password: u.Password,
	})
	if err != nil {
		//log.Println(err)
		return 0, err
	} else {
		return ID, nil
	}
}

func (userRepo *UserRepo) SelectByName(name string) *storage.UserEntity {
	userRepo.Lock()
	defer userRepo.Unlock()
	users := userRepo.Storage.FindAll()
	for _, value := range users {
		if value.UserName == name {
			return value
		}
	}
	return nil
}

func (userRepo *UserRepo) SelectById(uid uint) *storage.UserEntity {
	userRepo.Lock()
	defer userRepo.Unlock()
	uE, err := userRepo.Storage.Read(uid)
	if err != nil {
		println(err)
	}
	return uE
}

// Add the tweet id into the user
func (u *UserRepo) AddTweetToUser(uId uint, pId uint) bool {
	u.Lock()
	defer u.Unlock()
	userE, err := u.Storage.Read(uId)
	if err != nil {
		log.Println(err)
		return false
	} else {
		userE.Posts.PushBack(pId)
		u.Storage.Update(uId, userE)
		return true
	}
}

// Return all users in the database
func (u *UserRepo) FindAllUsers() []*storage.UserEntity {
	u.Lock()
	defer u.Unlock()
	ret := u.Storage.FindAll()
	return ret
}

// Check whether the user srcId follows the user targetId.
// Take O(#following) time
func (u *UserRepo) checkWhetherFollowing(srcId uint, targetId uint) bool {
	srcUserE, err := u.Storage.Read(srcId)
	if err != nil {
		log.Println(err)
		return false
	}
	if _, err := u.Storage.Read(targetId); err != nil {
		log.Println(err)
		return false
	}
	fl := srcUserE.Following
	for e := fl.Front(); e != nil; e = e.Next() {
		fuid := e.Value.(uint)
		if fuid == targetId {
			return true
		}
	}
	return false
}

func (u *UserRepo) CheckWhetherFollowing(srcId uint, targetId uint) bool {
	u.Lock()
	defer u.Unlock()
	return u.checkWhetherFollowing(srcId, targetId)
}

// User srcId starts to follow targetId.
func (u *UserRepo) StartFollowing(srcId uint, targetId uint) error {
	u.Lock()
	defer u.Unlock()
	srcUser, err1 := u.Storage.Read(srcId)
	targetUser, err2 := u.Storage.Read(targetId)
	if err1 == nil && err2 == nil {
		if srcUser.ID == targetUser.ID {
			return errors.New("cannot follow themselves")
		}
		if u.checkWhetherFollowing(srcId, targetId) {
			return errors.New("already followed")
		}
		srcUser.Following.PushBack(targetId) //This is wrong!
		targetUser.Follower.PushBack(srcId)
		u.Storage.Update(srcId, srcUser)
		u.Storage.Update(targetId, targetUser)
		return nil
	} else {
		log.Println(err1)
		log.Println(err2)
		return errors.New("illegal user id")
	}
}

// srcId stop following targetId.
// targetId remove the follower srcId.
func (u *UserRepo) StopFollowing(srcId uint, targetId uint) bool {
	u.Lock()
	defer u.Unlock()
	srcUser, err1 := u.Storage.Read(srcId)
	targetUser, err2 := u.Storage.Read(targetId)
	if err1 == nil && err2 == nil {
		if srcUser.ID == targetUser.ID {
			return false
		}
		for e := srcUser.Following.Front(); e != nil; e = e.Next() {
			v := e.Value.(uint)
			if v == targetId {
				srcUser.Following.Remove(e) // srcId stop following targetId
			}
		}
		for e := targetUser.Follower.Front(); e != nil; e = e.Next() {
			v := e.Value.(uint)
			if v == targetId {
				srcUser.Follower.Remove(e) // targetId remove the follower srcId
			}
		}
		u.Storage.Update(srcId, srcUser)
		u.Storage.Update(targetId, targetUser)
		return true
	} else {
		log.Println(err1)
		log.Println(err2)
		return false
	}
}
