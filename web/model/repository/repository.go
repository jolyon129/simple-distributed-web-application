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

type PostInfo struct {
	UserName string
	content  string
}

//TODO:
// Finish operations for repos!

func (userRepo *UserRepo) CreateNewUser(u *UserInfo) (uint, error) {
	ID, err := userRepo.Storage.Create(&storage.UserEntity{
		ID:       0,
		UserName: u.UserName,
		Password: u.Password,
	})
	if err != nil {
		log.Print(err)
		return 0, err
	} else {
		return ID, nil
	}
}

func (userRepo *UserRepo) SelectByName(name string) *storage.UserEntity {
	l := userRepo.Storage.FindAll()
	var next *list.Element
	for e := l.Front(); e != nil; e = next {
		u := e.Value.(storage.UserEntity)
		if u.UserName == name {
			return &u
		}
		next = e.Next()
	}
	return nil
}

func (userRepo *UserRepo) FindAll() *list.List {
	return userRepo.FindAll()
}

//func (UserRepo *UserRepo)
