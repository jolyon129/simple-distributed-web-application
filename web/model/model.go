package model

import (
	"zl2501-final-project/web/model/repository"
	"zl2501-final-project/web/model/storage"
	_ "zl2501-final-project/web/model/storage/memory"
)

var storageManager = storage.NewManager("memory")

func GetUserRepo() *repository.UserRepo {
	userRepo := repository.UserRepo{Storage: storageManager.UserStorage}
	return &userRepo
}
