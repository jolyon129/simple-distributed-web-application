package repository_test

import (
    "context"
    "fmt"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "log"
    "strconv"
    "sync"
    "time"
    "zl2501-final-project/backend/model"
    . "zl2501-final-project/backend/model/repository"
    "zl2501-final-project/backend/model/storage"
)

var userRepo *UserRepo
var tweetRepo *TweetRepo
var puId uint
var puId2 uint
var usersForTestFollowing []uint
var _ = BeforeSuite(func() {
    log.SetPrefix("Ginkgo: ")
    log.SetFlags(log.Ltime | log.Lshortfile)
    userRepo = model.GetUserRepo()
    tweetRepo = model.GetTweetRepo()
    //    result := make(chan uint)
    //    errorChan := make(chan error)
    timeout := 3000 * time.Millisecond
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    puId, _ = userRepo.CreateNewUser(ctx, &UserInfo{
        UserName: "jolyon129",
        Password: "123",
    })

    //    log.Print("puId1:", puId)
    puId2, _ = userRepo.CreateNewUser(ctx, &UserInfo{
        UserName: "jolyon2",
        Password: "123",
    })
    //    log.Print("puId2:", puId2)
    srcUser, _ := userRepo.SelectById(ctx, puId)
    usersForTestFollowing = make([]uint, 10)
    for i := 0; i < 10; i++ {
        id, _ := userRepo.CreateNewUser(ctx, &UserInfo{
            UserName: "userForTestFollowing" + strconv.Itoa(i),
            Password: "123",
        })
        userRepo.StartFollowing(ctx, srcUser.ID, id)
        usersForTestFollowing[i] = id
    }
    log.Print("Ready to Test")
})
var _ = Describe("User Repository", func() {
    timeout := 5 * time.Second
    Describe("Create New User in single thread", func() {
        Context("with a non-existed username", func() {
            It("should return a new User ID", func() {
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                uId, err := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "Zhuolun Li",
                    Password: "123",
                })
                _, _ = fmt.Fprintln(GinkgoWriter, "User ID: ", uId)
                Expect(err).Should(BeNil())
                Expect(uId).Should(Not(BeZero()))
                //                userE, _ := userRepo.storage.Read(uId)
                userE, _ := userRepo.SelectById(ctx, uId)
                Expect(userE.UserName).Should(Equal("Zhuolun Li"))
            })
        })
        Context("with a duplicated username", func() {
            It("should return error", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                fmt.Fprintln(GinkgoWriter, "The predefined User Id:", puId)
                uId, err := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "jolyon129",
                    Password: "123",
                })
                Expect(err).Should(Not(BeNil()))
                Expect(uId).Should(BeZero())
            })
        })
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                time.Sleep(500 * time.Millisecond)
                cancel()
                uId, err := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "Timeout",
                    Password: "123",
                })
                //log.Print(err)
                Expect(err).ShouldNot(BeNil())
                Expect(uId).Should(BeZero())
            })
        })
    })
    Describe("Create New User concurrently", func() {
        Context("with different user names", func() {
            It("should return different user entity", func() {
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                //                ctx := context.Background()
                var wg sync.WaitGroup
                num := 10
                wg.Add(num)
                uids := make([]uint, num)
                for i := 0; i < num; i++ {
                    go func(i int) {
                        defer GinkgoRecover()
                        defer wg.Done()
                        uid, err := userRepo.CreateNewUser(ctx, &UserInfo{
                            UserName: "user" + strconv.Itoa(i),
                            Password: "123",
                        })
                        if err != nil {
                            Fail(err.Error())
                        }
                        uids[i] = uid
                    }(i)
                }
                wg.Wait()
                for i := 0; i < num; i++ {
                    user, _ := userRepo.SelectById(ctx, uids[i])
                    Expect(user.UserName).Should(Equal(
                        "user" + strconv.Itoa(i)))
                }
            })
        })
        Context("with duplicated names", func() {
            It("should only succeed once", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                ch := make(chan int, 15)
                var wg sync.WaitGroup
                wg.Add(10)
                for i := 0; i < 10; i++ {
                    go func() {
                        defer wg.Done()
                        uid, err := userRepo.CreateNewUser(ctx, &UserInfo{
                            UserName: "dup",
                            Password: "123",
                        })
                        if err != nil && uid == 0 {
                            ch <- 1
                        }
                    }()
                }
                wg.Wait()
                Expect(len(ch)).Should(Equal(9))
            })
        })
    })
    Describe("Test SelectByName", func() {
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 1 * time.Second
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                time.Sleep(2 * time.Second)
                u, err := userRepo.SelectByName(ctx, "jolyon129")
                Expect(err).Should(Not(BeNil()))
                Expect(u).Should(BeZero())
            })
        })
        Context("with one existed name in two threads", func() {
            It("should return different user objects with same information correctly("+
                    "Different Copies)",
                func() {
                    ctx, cancel := context.WithTimeout(context.Background(), timeout)
                    defer cancel()
                    ch := make(chan *storage.UserEntity, 2)
                    for i := 0; i < 2; i++ {
                        go func() {
                            userEntity, _ := userRepo.SelectByName(ctx, "jolyon129")
                            ch <- userEntity
                        }()
                    }
                    u1 := <-ch
                    u2 := <-ch
                    Expect(u1).Should(Not(BeIdenticalTo(u2)))
                    Expect(u1.UserName).Should(Equal(u2.UserName))
                })
        })
        Context("with non-existed name", func() {
            It("should return nil", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                u, _ := userRepo.SelectByName(ctx, "fake name")
                Expect(u).Should(BeNil())
            })
        })
    })
    Describe("Test SelectById", func() {
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                time.Sleep(1000 * time.Millisecond)
                u, err := userRepo.SelectById(ctx, puId)
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("with one existed id in two threads", func() {
            It("should return different pointers with same information(Two Copies)", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                //                log.Print("Im here")
                ch := make(chan *storage.UserEntity, 2)
                for i := 0; i < 2; i++ {
                    go func() {
                        defer GinkgoRecover()
                        u, err := userRepo.SelectById(ctx, puId)
                        ch <- u
                        if err != nil {
                            Fail(err.Error())
                        }
                    }()
                }
                u1 := <-ch
                u2 := <-ch
                Expect(u1).Should(Not(BeIdenticalTo(u2)))
                Expect(u1.Following).Should(Not(BeIdenticalTo(u2.Following)))
                Expect(u1.Following.Front()).Should(Not(BeIdenticalTo(u2.Following.Front())))
                Expect(u1.UserName).Should(Equal(u2.UserName))
            })
        })
        Context("with non-existed name", func() {
            It("should return nil", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                var uid uint
                uid = 222222
                u, _ := userRepo.SelectById(ctx, uid)
                Expect(u).Should(BeNil())
            })
        })
    })
    Describe("Add Tweets To User", func() {
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                time.Sleep(1000 * time.Millisecond)
                u, err := userRepo.AddTweetToUser(ctx, puId, uint(10))
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("When create tweets and add tweet concurrently", func() {
            It("Should be synchronized", func() {
                //ctx, cancel := context.WithTimeout(context.Background(), timeout)
                //defer cancel()
                //nuId, _ :=userRepo.CreateNewUser(ctx, &UserInfo{
                //    UserName: "newuser2",
                //    Password: "123",
                //})

            })
        })
        Context("When create and save 10 tweets to 2 users concurrently", func() {
            It("should be synchronized", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                var wg sync.WaitGroup
                wg.Add(20)
                nuId, _ := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "newuser2",
                    Password: "123",
                })
                nuId2, _ := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "newuser3",
                    Password: "123",
                })

                for i := 0; i < 10; i++ {
                    go func(i int) {
                        defer wg.Done()
                        tId, _ := tweetRepo.SaveTweet(ctx, TweetInfo{
                            UserID:  nuId2,
                            Content: strconv.Itoa(100 + i),
                        })
                        userRepo.AddTweetToUser(ctx, nuId2, tId)
                    }(i)
                    go func(i int) {
                        defer wg.Done()
                        tId, _ := tweetRepo.SaveTweet(ctx, TweetInfo{
                            UserID:  nuId,
                            Content: strconv.Itoa(10 + i),
                        })
                        userRepo.AddTweetToUser(ctx, nuId, tId)
                    }(i)
                }
                wg.Wait()
                u, _ := userRepo.SelectById(ctx, nuId2)
                checkSum2 := 0
                for e := u.Tweets.Front(); e != nil; e = e.Next() {
                    tId := e.Value.(uint)
                    costr, _ := tweetRepo.SelectById(ctx, tId)
                    coint, _ := strconv.Atoi(costr.Content)
                    checkSum2 += coint
                }
                Expect(checkSum2).Should(Equal(1045))

                checkSUm1 := 0
                u2, _ := userRepo.SelectById(ctx, nuId)
                for e := u2.Tweets.Front(); e != nil; e = e.Next() {
                    tId := e.Value.(uint)
                    costr, _ := tweetRepo.SelectById(ctx, tId)
                    coint, _ := strconv.Atoi(costr.Content)
                    checkSUm1 += coint
                }
                Expect(checkSUm1).Should(Equal(145))
            })
        })
        Context("Add 10 tweets concurrently to the same user", func() {
            It("should create all of the tweets correctly", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                var wg sync.WaitGroup
                num := 1
                wg.Add(num)
                uId, err := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "whatsurpronlem",
                    Password: "123",
                })
                if err != nil {
                    log.Print(uId)
                    Fail(err.Error())
                }
                for i := 0; i < num; i++ {
                    go func(i int) {
                        defer wg.Done()
                        defer GinkgoRecover()
                        _, err := userRepo.AddTweetToUser(ctx, uId, uint(i))
                        if err != nil {
                            Fail(err.Error())
                        }
                    }(i)
                }
                wg.Wait()
                uE, err2 := userRepo.SelectById(ctx, uId)
                if err2 != nil {
                    Fail(err2.Error())
                }
                Expect(uE.Tweets.Len()).Should(Equal(num))
            })
        })
    })
    Describe("Find All Users", func() {
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                time.Sleep(1000 * time.Millisecond)
                u, err := userRepo.FindAllUsers(ctx)
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("When read concurrently", func() {
            It("should return different pointers with same info", func() {
                var wg sync.WaitGroup
                wg.Add(10)
                pointerArr := make([][]*storage.UserEntity, 10)
                for i := 0; i < 10; i++ {
                    go func(i int) {
                        ctx, cancel := context.WithTimeout(context.Background(), timeout)
                        defer cancel()
                        defer wg.Done()
                        pointerArr[i], _ = userRepo.FindAllUsers(ctx)
                    }(i)
                }
                wg.Wait()
                for i := 0; i < 9; i++ {
                    Expect(pointerArr[i]).Should(Not(BeIdenticalTo(pointerArr[i+1])))
                    Expect(len(pointerArr[i])).Should(Equal(len(pointerArr[i+1])))
                    for j := 0; j < len(pointerArr[i]); j++ {
                        Expect(pointerArr[i][j].ID).Should(Equal(pointerArr[i+1][j].ID))
                    }
                }
            })
        })
    })
    Describe("Check Whether following", func() {
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                time.Sleep(1000 * time.Millisecond)
                u, err := userRepo.CheckWhetherFollowing(ctx, puId, puId2)
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("when check one whom I already followed ", func() {
            It("should return true", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                uId1, _ := userRepo.CreateNewUser(ctx, &UserInfo{
                    UserName: "testfollowing",
                    Password: "123",
                })
                userRepo.StartFollowing(ctx, puId, uId1)
                r, _ := userRepo.CheckWhetherFollowing(ctx, puId, uId1)
                Expect(r).Should(BeTrue())
            })
        })
        Context("when check a lot whom I was already following concurrently", func() {
            It("all return true", func() {
                //TODO:
                //	Check why this failed?
                for i := 0; i < 10; i++ {
                    go func(t int) {
                        ctx, cancel := context.WithTimeout(context.Background(), timeout)
                        defer cancel()
                        defer GinkgoRecover()
                        res, _ := userRepo.CheckWhetherFollowing(ctx, puId, usersForTestFollowing[t])
                        Expect(res).Should(BeTrue())
                    }(i)
                }
            })
        })
    })
    Describe("Start/Stop following", func() {
        Context("when start following with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                time.Sleep(1000 * time.Millisecond)
                u, err := userRepo.StartFollowing(ctx, puId, puId2)
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("when stop following with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                time.Sleep(1000 * time.Millisecond)
                cancel()
                u, err := userRepo.StopFollowing(ctx, puId, puId2)
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("When start following a lot of people concurrently ", func() {
            It("should follow all of them", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                srcUser, _ := userRepo.SelectById(ctx, puId)
                users := make([]uint, 10)
                var wg1 sync.WaitGroup
                wg1.Add(10)
                //log.Println("Before Length:", srcUser.Following.Len())
                for i := 0; i < 10; i++ {
                    go func(t int) {
                        defer wg1.Done()
                        defer GinkgoRecover()
                        id, err0 := userRepo.CreateNewUser(ctx, &UserInfo{
                            UserName: "newUserForTestFollowing" + strconv.Itoa(t),
                            Password: "123",
                        })
                        Expect(err0).Should(BeNil())
                        _, err := userRepo.StartFollowing(ctx, srcUser.ID, id)
                        Expect(err).Should(BeNil())
                        //log.Println("Checking again", id)
                        users[t] = id
                    }(i)
                }
                wg1.Wait()
                ctx1, cancel1 := context.WithTimeout(context.Background(), timeout)
                defer cancel1()
                //log.Println("Finished all waits")
                for i := 0; i < len(users); i++ {
                    //log.Println(users[i])
                    idToTest := users[i]
                    res, err := userRepo.CheckWhetherFollowing(ctx1, srcUser.ID, idToTest)
                    if err != nil {
                        Fail(err.Error())
                    }
                    Expect(res).Should(BeTrue())
                }
            })
        })
        Context("When stop following concurrently", func() {
            It("should stop following", func() {
                ctx, _ := context.WithTimeout(context.Background(), timeout)
                //                defer cancel()
                for i := 0; i < 5; i++ {
                    go func(i int) {
                        defer GinkgoRecover()
                        userRepo.StopFollowing(ctx, puId, usersForTestFollowing[i])
                        res, err := userRepo.CheckWhetherFollowing(ctx, puId,
                            usersForTestFollowing[i])
                        if err != nil {
                            Fail(err.Error())
                        }
                        Expect(res).Should(Not(BeTrue()))
                    }(i)
                }
            })
        })
    })

})
var _ = Describe("Tweet Repository", func() {
    timeout := 5 * time.Second
    Describe("Create new post/Tweet", func() {
        Context("with a timeout context", func() {
            It("should return timeout err and wipe out the redundant tweet", func() {
                timeout := 200 * time.Millisecond
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                time.Sleep(500 * time.Millisecond)
                cancel()
                u, err := tweetRepo.SaveTweet(ctx, TweetInfo{})
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("When tweet one message", func() {
            It("should succeed", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                pid, err := tweetRepo.SaveTweet(ctx, TweetInfo{
                    UserID:  puId,
                    Content: "Test",
                })
                if err != nil {
                    Fail(err.Error())
                }
                Expect(pid).To(Not(BeZero()))
                //pE, _ := tweetRepo.storage.Read(pid)
                pE, _ := tweetRepo.SelectById(ctx, pid)
                Expect(pE.UserID).To(Equal(puId))
            })
        })
        Context("When post 20 tweets concurrently", func() {
            It("should store 20 tweets correctly", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                var wg sync.WaitGroup
                wg.Add(20)
                postIds := make([]uint, 20)
                for i := 0; i < 20; i++ {
                    go func(i int) {
                        defer wg.Done()
                        pid, err := tweetRepo.SaveTweet(ctx, TweetInfo{
                            UserID:  puId,
                            Content: "TestConcurrency" + strconv.Itoa(i),
                        })
                        Expect(err).To(BeNil())
                        Expect(pid).To(Not(BeZero()))
                        postIds[i] = pid
                    }(i)
                }
                wg.Wait()
                for i := 0; i < 20; i++ {
                    //pE, err := tweetRepo.storage.Read(postIds[i])
                    pE, err := tweetRepo.SelectById(ctx, postIds[i])
                    Expect(err).Should(BeNil())
                    Expect(pE.Content).Should(Equal("TestConcurrency" + strconv.Itoa(i)))
                }
            })
        })
    })
    Describe("Select Tweet By Id", func() {
        Context("with a timeout context", func() {
            It("should return timeout err", func() {
                timeout := 200 * time.Millisecond
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                time.Sleep(500 * time.Millisecond)
                cancel()
                u, err := tweetRepo.SelectById(ctx, uint(101))
                Expect(err).ShouldNot(BeNil())
                Expect(u).Should(BeZero())
            })
        })
        Context("Read the same tweet 20 times concurrently", func() {
            It("should return 20 different pointers with the same information", func() {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                defer cancel()
                var wg sync.WaitGroup
                wg.Add(20)
                postIds := make([]uint, 20)
                for i := 0; i < 20; i++ {
                    go func(i int) {
                        defer wg.Done()
                        defer GinkgoRecover()
                        pid, err := tweetRepo.SaveTweet(ctx, TweetInfo{
                            UserID:  puId,
                            Content: "TestConcurrency",
                        })
                        Expect(err).To(BeNil())
                        Expect(pid).To(Not(BeZero()))
                        postIds[i] = pid
                    }(i)
                }
                wg.Wait()
                postEs := make([]*storage.TweetEntity, 20)
                wg.Add(20)
                for i := 0; i < 20; i++ {
                    go func(i int) {
                        defer wg.Done()
                        post, _ := tweetRepo.SelectById(ctx, postIds[i])
                        postEs[i] = post
                    }(i)
                }
                wg.Wait()
                prev, _ := tweetRepo.SelectById(ctx, postIds[0])
                Expect(prev.Content).Should(Equal("TestConcurrency"))
                for i := 1; i < 20; i++ {
                    post := postEs[i]
                    Expect(post).ShouldNot(BeIdenticalTo(prev))
                    Expect(post.Content).Should(Equal(prev.Content))
                    prev = post
                }
            })
        })
    })
})
