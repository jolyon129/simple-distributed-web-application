package main

import (
    "context"
    "flag"
    "fmt"
    "math/rand"
    "time"
    "zl2501-final-project/backend"
    "zl2501-final-project/backend/model"
    "zl2501-final-project/backend/model/repository"
)

func main() {
    fmt.Println("This is the backend service")
    addData := flag.Bool("data", false, "add default data")
    flag.Parse()
    addDefaultData(*addData)
    backend.StartService()
}

func addDefaultData(addData bool) {
    if !addData {
        return
    }
    bgctx := context.Background()
    userRepo := model.GetUserRepo()
    hash, _ := backend.EncodePassword("123")
    rand.Seed(time.Now().UnixNano())
    uId, _ := userRepo.CreateNewUser(bgctx, &repository.UserInfo{
        UserName: "zl2501",
        Password: hash,
    })
    tweetRepo := model.GetTweetRepo()
    pid1, _ := tweetRepo.SaveTweet(bgctx, repository.TweetInfo{
        UserID:  uId,
        Content: "This is my first tweet!",
    })
    userRepo.AddTweetToUser(bgctx, uId, pid1)
    time.Sleep(2 * time.Second)
    pid2, _ := tweetRepo.SaveTweet(bgctx, repository.TweetInfo{
        UserID:  uId,
        Content: "I really hope this coronavirus is over very soon! No more quarantine!",
    })
    userRepo.AddTweetToUser(bgctx, uId, pid2)

    uId2, _ := userRepo.CreateNewUser(bgctx, &repository.UserInfo{
        UserName: "jolyon129",
        Password: hash,
    })
    pid3, _ := tweetRepo.SaveTweet(bgctx, repository.TweetInfo{
        UserID:  uId2,
        Content: "Gotta give him a new blanket when I back home. #NationalPuppyDay #Westie",
    })
    userRepo.AddTweetToUser(bgctx, uId2, pid3)
    time.Sleep(1 * time.Second)
    pid4, _ := tweetRepo.SaveTweet(bgctx, repository.TweetInfo{
        UserID:  uId2,
        Content: "BTW, this is Sakuragi. 3 year’s old. And his name is from a Japanese anime. ",
    })
    userRepo.AddTweetToUser(bgctx, uId2, pid4)
    userRepo.StartFollowing(bgctx, uId, uId2)
}
