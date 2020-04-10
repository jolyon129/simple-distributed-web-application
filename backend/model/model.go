package model

import (
	"sync"
	"zl2501-final-project/backend/model/repository"
	"zl2501-final-project/backend/model/storage"
	_ "zl2501-final-project/backend/model/storage/memory"
)

var storageManager = storage.NewManager("memory")
var userRepo *repository.UserRepo
var postRepo *repository.PostRepo
var muForUser sync.Mutex
var muForPost sync.Mutex

// Get the singleton of user repository
// This is synchronized bc multiple threads can call this at the same time
func GetUserRepo() *repository.UserRepo {
	muForUser.Lock()
	defer muForUser.Unlock()
	if userRepo == nil {
		userRepo = repository.NewUserRepo()
		userRepo.Storage = storageManager.UserStorage
		return userRepo
	} else {
		return userRepo
	}
}

// Get the singleton of post repository
// This is synchronized bc multiple threads can call this at the same time
func GetPostRepo() *repository.PostRepo {
	muForPost.Lock()
	defer muForPost.Unlock()
	if postRepo == nil {
		postRepo = &repository.PostRepo{Storage: storageManager.TweetStorage}
		return postRepo
	} else {
		return postRepo
	}
}
