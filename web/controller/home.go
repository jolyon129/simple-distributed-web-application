package controller

import (
    "context"
    "html/template"
    "log"
    "net/http"
    "time"
    "zl2501-final-project/web/constant"
    . "zl2501-final-project/web/pb"
)

// Cache the result from tweetId to Tweet
var tweetCacheMap = make(map[uint64]*TweetEntity)

// Cache the result from userId to User
var userCacheMap = make(map[uint64]*UserEntity)

type tweet struct {
    Content   string
    CreatedAt string
    CreatedBy string
    UserId    int
}

type homeView struct {
    Name     string
    MyTweets []tweet
    Feed     []tweet
}

func Home(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        userId, err1 := GetMyUserId(r)
        if err1!=nil{
            log.Printf(err1.Error())
        }
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        res, _ := BackendClientIns.UserSelectById(ctx, &SelectByIdRequest{
            Id: userId,
        })

        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "home.html")
        w.Header().Set("Content-Type", "text/html")
        view := homeView{Name: res.User.UserName, MyTweets: make([]tweet, 0)}
        userE := res.User
        userCacheMap[userE.UserId] = userE

        // Iterate in reverse order because the latest one is stored in the tail in DB
        tweets := make([]tweet, 0)
        for i := len(userE.Tweets) - 1; i >= 0; i-- {
            ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
            res, _ := BackendClientIns.TweetSelectById(ctx, &SelectByIdRequest{
                Id: userE.Tweets[i],
            })
            tweets = append(tweets, tweet{
                Content:   res.Msg.Content,
                CreatedAt: res.Msg.CreatedTime,
                CreatedBy: userE.UserName,
                UserId:    int(userE.UserId),
            })
        }

        // Build feed
        followingUIdList := userE.Followings
        tweetsInOldestFirst := make([]uint64, 0)
        //for e := followingUIdList.Front(); e != nil; e = e.Next() {
        //    u := userRepo.SelectById(e.Value.(uint))
        //    tweetsInOldestFirst = mergeSortedPostList(tweetsInOldestFirst, u.Posts) // Keep merging into the new feed
        //}
        for _, uId := range followingUIdList {
            ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
            res, _ := BackendClientIns.UserSelectById(ctx, &SelectByIdRequest{
                Id: uId,
            })
            userCacheMap[res.User.UserId] = res.User
            tweetsInOldestFirst = mergeSortedTweets(tweetsInOldestFirst, res.User.Tweets)
        }

        tweetsInOldestFirst = mergeSortedTweets(tweetsInOldestFirst, userE.Tweets) // Merge my own tweets
        retFeed := make([]tweet, len(tweetsInOldestFirst))
        z := 0
        // Iterate in revers oder so that the latest comes first
        for i := len(tweetsInOldestFirst) - 1; i >= 0; i-- {
            t := tweetCacheMap[tweetsInOldestFirst[i]]
            retFeed[z] = tweet{
                Content:   t.Content,
                CreatedAt: t.CreatedTime,
                CreatedBy: userCacheMap[t.UserId].UserName,
                UserId:    int(userCacheMap[t.UserId].UserId),
            }
            z++
        }
        log.Println(retFeed)
        view.Feed = retFeed
        t.Execute(w, view)
    }
}

// Return a list of tweet (from oldest to newest)
func mergeSortedTweets(l1 []uint64, l2 []uint64) []uint64 {
    ret := make([]uint64, len(l1)+len(l2))
    i := 0
    j := 0
    z := 0
    for i != len(l1) && j != len(l2) {
        tid1 := l1[i]
        tid2 := l2[j]
        tweet1 := getTweetById(tid1)
        tweet2 := getTweetById(tid2)
        ts1, _ := time.Parse(constant.TimeFormat, tweet1.Msg.CreatedTime)
        ts2, _ := time.Parse(constant.TimeFormat, tweet2.Msg.CreatedTime)
        tweetCacheMap[tid1] = tweet1.Msg
        tweetCacheMap[tid2] = tweet2.Msg
        if ts1.Before(ts2) {
            ret[z] = tid1
            i++
        } else {
            ret[z] = tid2
            j++
        }
        z++
    }
    for i < len(l1) {
        tid1 := l1[i]
        tweet1 := getTweetById(tid1)
        tweetCacheMap[tid1] = tweet1.Msg
        ret[z] = tid1
        i++
        z++
    }
    for j < len(l2) {
        tid2 := l2[j]
        tweet2 := getTweetById(tid2)
        tweetCacheMap[tid2] = tweet2.Msg
        ret[z] = tid2
        j++
        z++
    }
    return ret
}

func getTweetById(id uint64) *TweetSelectByIdResponse {
    ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
    res1, _ := BackendClientIns.TweetSelectById(ctx, &SelectByIdRequest{
        Id: id,
    })
    return res1
}
