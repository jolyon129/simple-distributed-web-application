package model

import (
	"zl2501-final-project/web/model/repository"
	"zl2501-final-project/web/model/storage"
	_ "zl2501-final-project/web/model/storage/memory"
)

var storageManager = storage.NewManager("memory")
var userRepo *repository.UserRepo

// Get the singleton of user repository
func GetUserRepo() *repository.UserRepo {
	if userRepo == nil {
		userRepo = &repository.UserRepo{Storage: storageManager.UserStorage}
		return userRepo
	} else {
		return userRepo
	}
}
