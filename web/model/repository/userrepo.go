package repository

import (
	"container/list"
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

func (userRepo *UserRepo) CreateNewUser(u *UserInfo) (uint, error) {
	ID, err := userRepo.Storage.Create(&storage.UserEntity{
		ID:       0,
		UserName: u.UserName,
		Password: u.Password,
	})
	if err != nil {
		log.Println(err)
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
	uE, err:=userRepo.Storage.Read(uid)
	if err!=nil{
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
		u.Storage.Update(uId, &storage.UserEntity{
			Password: userE.Password,
			Posts:    userE.Posts,
		})
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
func (u *UserRepo) CheckWhetherFollow(srcId uint, targetId uint) bool {
	srcUserE, err := u.Storage.Read(srcId)
	if err != nil {
		log.Println(err)
		return false
	}
	if _, err := u.Storage.Read(targetId); err!=nil{
		return false
	}
	fl := srcUserE.Follower
	for e := fl.Front(); e != nil; e.Next() {
		fuid := e.Value.(uint)
		if fuid == targetId {
			return true
		}
	}
	return false

}
