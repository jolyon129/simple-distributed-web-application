package main

import (
	"os"
	"zl2501-final-project/web"
	"zl2501-final-project/web/controller"
	"zl2501-final-project/web/model"
	"zl2501-final-project/web/model/repository"
)

func main() {
	// Change working directory to web/ first
	err := os.Chdir("../../web")
	if err != nil {
		panic(err)
	}
	addPredefinedUsers()
	web.StartService()
}

func addPredefinedUsers() {
	userRepo := model.GetUserRepo()
	hash, _ := controller.EncodePassword("123")
	userRepo.CreateNewUser(&repository.UserInfo{
		UserName: "zl2501",
		Password: hash,
	})
}
