package model

import (
	"zl2501-final-project/web/model/repository"
	"zl2501-final-project/web/model/storage"
	_ "zl2501-final-project/web/model/storage/memory"
)

var storageManager = storage.NewManager("memory")
var userRepo *repository.UserRepo
var postRepo *repository.PostRepo

// Get the singleton of user repository
func GetUserRepo() *repository.UserRepo {
	if userRepo == nil {
		userRepo = &repository.UserRepo{Storage: storageManager.UserStorage}
		return userRepo
	} else {
		return userRepo
	}
}

// Get the singleton of post repository
func GetPostRepo() *repository.PostRepo {
	if postRepo == nil {
		postRepo = &repository.PostRepo{Storage: storageManager.PostStorage}
		return postRepo
	} else {
		return postRepo
	}
}
