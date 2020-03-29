package repository

import (
	"container/list"
	"errors"
	"log"
	"zl2501-final-project/web/model/storage"
)

type UserRepo struct {
	Storage storage.UserStorageInterface
}

type UserInfo struct {
	UserName string
	Password string
}

// Create a new user and return user id
// if the user name is no duplicated.
// Otherwise return error
func (userRepo *UserRepo) CreateNewUser(u *UserInfo) (uint, error) {
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
	l := userRepo.Storage.FindAll()
	var next *list.Element
	for e := l.Front(); e != nil; e = next {
		u := e.Value.(*storage.UserEntity)
		if u.UserName == name {
			return u
		}
		next = e.Next()
	}
	return nil
}

func (userRepo *UserRepo) SelectById(uid uint) *storage.UserEntity {
	uE, err := userRepo.Storage.Read(uid)
	if err != nil {
		println(err)
	}
	return uE
}

// Add the tweet id into the user
func (u *UserRepo) AddTweetToUser(uId uint, pId uint) bool {
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
	l := u.Storage.FindAll()
	ret := make([]*storage.UserEntity, 0)
	for e := l.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*storage.UserEntity))
	}
	return ret
}

// Check whether the user srcId follows the user targetId.
// Take O(#following) time
func (u *UserRepo) CheckWhetherFollowing(srcId uint, targetId uint) bool {
	srcUserE, err := u.Storage.Read(srcId)
	if err != nil {
		log.Println(err)
		return false
	}
	if _, err := u.Storage.Read(targetId); err != nil {
		return false
	}
	fl := srcUserE.Following
	for e := fl.Front(); e != nil; e=e.Next() {
		fuid := e.Value.(uint)
		if fuid == targetId {
			return true
		}
	}
	return false
}

// User srcId starts to follow targetId.
func (u *UserRepo) StartFollowing(srcId uint, targetId uint) error {
	srcUser, err1 := u.Storage.Read(srcId)
	targetUser, err2 := u.Storage.Read(targetId)
	if err1 == nil && err2 == nil {
		if srcUser.ID == targetUser.ID {
			return errors.New("cannot follow themselves")
		}
		if u.CheckWhetherFollowing(srcId,targetId){
			return errors.New("already followed")
		}
		srcUser.Following.PushBack(targetId)
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
