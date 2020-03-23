package main

import (
	"os"
	"time"
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
	addDefaultData()
	web.StartService()
}

func addDefaultData() {
	userRepo := model.GetUserRepo()
	hash, _ := controller.EncodePassword("123")
	uId, _ := userRepo.CreateNewUser(&repository.UserInfo{
		UserName: "zl2501",
		Password: hash,
	})
	postRepo := model.GetPostRepo()
	pid1, _ := postRepo.CreateNewPost(repository.PostInfo{
		UserID:  uId,
		Content: "This is my first tweet!",
	})
	userRepo.AddTweetToUser(uId, pid1)
	time.Sleep(2 * time.Second)
	pid2, _ := postRepo.CreateNewPost(repository.PostInfo{
		UserID:  uId,
		Content: "I really hope this coronavirus can end soon! No more quarantine!",
	})
	userRepo.AddTweetToUser(uId, pid2)

	uId2, _ := userRepo.CreateNewUser(&repository.UserInfo{
		UserName: "jolyon129",
		Password: hash,
	})
	pid3, _ := postRepo.CreateNewPost(repository.PostInfo{
		UserID:  uId2,
		Content: "Gotta give him a new blanket when I back home. #NationalPuppyDay #Westie",
	})
	userRepo.AddTweetToUser(uId2, pid3)
	time.Sleep(2 * time.Second)
	pid4, _ := postRepo.CreateNewPost(repository.PostInfo{
		UserID:  uId2,
		Content: "BTW, this is Sakuragi. 3 year’s old. And his name is from a Japanese anime. ",
	})
	userRepo.AddTweetToUser(uId2, pid4)
}
