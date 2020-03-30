package repository_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"strconv"
	"sync"
	"zl2501-final-project/web/model"
	. "zl2501-final-project/web/model/repository"
	"zl2501-final-project/web/model/storage"
)

var userRepo *UserRepo
var postRepo *PostRepo
var puId uint
var puId2 uint
var usersForTestFollowing []uint
var _ = BeforeSuite(func() {
	userRepo = model.GetUserRepo()
	postRepo = model.GetPostRepo()
	puId, _ = userRepo.CreateNewUser(&UserInfo{
		UserName: "jolyon129",
		Password: "123",
	})
	puId2, _ = userRepo.CreateNewUser(&UserInfo{
		UserName: "jolyon2",
		Password: "123",
	})
	srcUser := userRepo.SelectById(puId)
	usersForTestFollowing = make([]uint, 10)
	for i := 0; i < 10; i++ {
		id, _ := userRepo.CreateNewUser(&UserInfo{
			UserName: "userForTestFollowing" + strconv.Itoa(i),
			Password: "123",
		})
		userRepo.StartFollowing(srcUser.ID, id)
		usersForTestFollowing[i] = id
	}
})
var _ = Describe("User Repository", func() {
	Describe("Create New User in single thread", func() {
		Context("with a non-existed username", func() {
			It("should return a new User ID", func() {
				uId, err := userRepo.CreateNewUser(&UserInfo{
					UserName: "Zhuolun Li",
					Password: "123",
				})
				_, _ = fmt.Fprintln(GinkgoWriter, "User ID: ", uId)
				Expect(err).Should(BeNil())
				Expect(uId).Should(Not(BeZero()))
				userE, _ := userRepo.Storage.Read(uId)
				Expect(userE.UserName).Should(Equal("Zhuolun Li"))
			})
		})
		Context("with a duplicated username", func() {
			It("should return error", func() {

				fmt.Fprintln(GinkgoWriter, "The predefined User Id:", puId)
				uId, err := userRepo.CreateNewUser(&UserInfo{
					UserName: "jolyon129",
					Password: "123",
				})
				Expect(err).Should(Not(BeNil()))
				Expect(uId).Should(BeZero())
			})
		})
	})
	Describe("Create New User concurrently", func() {
		Context("with different user names", func() {
			It("should return different user entity", func() {
				var wg sync.WaitGroup
				wg.Add(10)
				uids := make([]uint, 10)
				for i := 0; i < 10; i++ {
					go func(i int) {
						defer wg.Done()
						uid, _ := userRepo.CreateNewUser(&UserInfo{
							UserName: "user" + strconv.Itoa(i),
							Password: "123",
						})
						uids[i] = uid
					}(i)
				}
				wg.Wait()
				for i := 0; i < 10; i++ {
					Expect(userRepo.SelectById(uids[i]).UserName).Should(Equal("user" + strconv.Itoa(i)))
				}
			})
		})
		Context("with duplicated names", func() {
			It("should only succeed once", func() {
				ch := make(chan int, 15)
				var wg sync.WaitGroup
				wg.Add(10)
				for i := 0; i < 10; i++ {
					go func() {
						defer wg.Done()
						uid, err := userRepo.CreateNewUser(&UserInfo{
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
		Context("with one existed name in two threads", func() {
			It("should return different user objects with same information correctly", func() {
				ch := make(chan *storage.UserEntity, 2)
				for i := 0; i < 2; i++ {
					go func() {
						userEntity := userRepo.SelectByName("jolyon129")
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
				u := userRepo.SelectByName("fake name")
				Expect(u).Should(BeNil())
			})
		})
	})
	Describe("Test SelectById", func() {
		Context("with one existed id in two threads", func() {
			It("should return different pointers with same information", func() {
				ch := make(chan *storage.UserEntity)
				for i := 0; i < 2; i++ {
					go func() {
						ch <- userRepo.SelectById(puId)
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
				var uid uint
				uid = 222222
				u := userRepo.SelectById(uid)
				Expect(u).Should(BeNil())
			})
		})
	})
	Describe("Add Tweet To User", func() {
		Context("When Adding to users concurrently", func() {
			It("should be synchronized", func() {
				var wg sync.WaitGroup
				wg.Add(20)
				nuId, _ := userRepo.CreateNewUser(&UserInfo{
					UserName: "newuser",
					Password: "123",
				})
				for i := 0; i < 10; i++ {
					go func(i int) {
						defer wg.Done()
						userRepo.AddTweetToUser(puId, uint(100+i))
					}(i)
					go func(i int) {
						defer wg.Done()
						userRepo.AddTweetToUser(nuId, uint(10+i))
					}(i)
				}
				wg.Wait()
				u := userRepo.SelectById(puId)
				checkSum := uint(0)
				for e := u.Posts.Front(); e != nil; e = e.Next() {
					pId := e.Value.(uint)
					checkSum += pId
				}
				Expect(checkSum).Should(Equal(uint(1045)))

				checkSum2 := uint(0)
				u2 := userRepo.SelectById(nuId)
				for e := u2.Posts.Front(); e != nil; e = e.Next() {
					pId := e.Value.(uint)
					checkSum2 += pId
				}
				Expect(checkSum2).Should(Equal(uint(145)))
			})
		})

	})
	Describe("Find All Users", func() {
		Context("When read concurrently", func() {
			It("should return different pointers with same info", func() {
				var wg sync.WaitGroup
				wg.Add(10)
				pointerArr := make([][]*storage.UserEntity, 10)
				for i := 0; i < 10; i++ {
					go func(i int) {
						defer wg.Done()
						pointerArr[i] = userRepo.FindAllUsers()
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
		Context("when check one whom I already followed ", func() {
			It("should return true", func() {
				uId1, _ := userRepo.CreateNewUser(&UserInfo{
					UserName: "testfollowing",
					Password: "123",
				})
				userRepo.StartFollowing(puId, uId1)
				Expect(userRepo.CheckWhetherFollowing(puId, uId1)).Should(BeTrue())
			})
		})
		Context("when check a lot whom I already followed concurrently", func() {
			It("all return true", func() {
				for i := 0; i < 10; i++ {
					go func(t int) {
						res := userRepo.CheckWhetherFollowing(puId, usersForTestFollowing[t])
						Expect(res).Should(BeTrue())
					}(i)
				}
			})
		})
	})
	Describe("Start/Stop following", func() {
		Context("When start following a lot of people concurrently ", func() {
			It("should follow all of them", func() {
				srcUser := userRepo.SelectById(puId)
				users := make([]uint, 10)
				var wg1 sync.WaitGroup
				wg1.Add(10)
				//log.Println("Before Length:", srcUser.Following.Len())
				for i := 0; i < 10; i++ {
					go func(t int) {
						defer wg1.Done()
						defer GinkgoRecover()
						id, err0 := userRepo.CreateNewUser(&UserInfo{
							UserName: "newUserForTestFollowing" + strconv.Itoa(t),
							Password: "123",
						})
						Expect(err0).Should(BeNil())
						err := userRepo.StartFollowing(srcUser.ID, id)
						Expect(err).Should(BeNil())
						//log.Println("Checking again", id)
						users[t] = id
					}(i)
				}
				wg1.Wait()
				//log.Println("Finished all waits")
				myuser := userRepo.SelectById(puId)
				log.Println("Length:", myuser.Following.Len())
				for i := 0; i < len(users); i++ {
					//log.Println(users[i])
					idToTest := users[i]
					res := userRepo.CheckWhetherFollowing(srcUser.ID, idToTest)
					Expect(res).Should(BeTrue())
				}
			})
		})
		Context("When stop following concurrently", func() {
			It("should stop following", func() {
				for i := 0; i < 5; i++ {
					go func(i int) {
						defer GinkgoRecover()
						userRepo.StopFollowing(puId, usersForTestFollowing[i])
						Expect(userRepo.CheckWhetherFollowing(puId, usersForTestFollowing[i])).Should(Not(BeTrue()))
					}(i)
				}
			})
		})
	})
	Describe("Add Tweets", func() {
		Context("Add 10 tweets concurrently", func() {
			It("should create all of the tweets correctly", func() {
				var wg sync.WaitGroup
				wg.Add(10)
				uId, _ := userRepo.CreateNewUser(&UserInfo{
					UserName: "TestForTweeting",
					Password: "123",
				})
				for i := 0; i < 10; i++ {
					go func(i int) {
						defer wg.Done()
						userRepo.AddTweetToUser(uId, uint(i))
					}(i)
				}
				wg.Wait()
				uE := userRepo.SelectById(uId)
				Expect(uE.Posts.Len()).Should(Equal(10))
			})
		})
	})

})
var _ = Describe("Post Repository", func() {
	Describe("Create new post/Tweet", func() {
		Context("When tweet one message", func() {
			It("should succeed", func() {
				pid, err := postRepo.CreateNewPost(PostInfo{
					UserID:  puId,
					Content: "Test",
				})
				Expect(err).To(BeNil())
				Expect(pid).To(Not(BeZero()))
				pE, _ := postRepo.Storage.Read(pid)
				Expect(pE.UserID).To(Equal(puId))
			})
		})
		Context("When post 20 tweets concurrently", func() {
			It("should store 20 tweets correctly", func() {
				var wg sync.WaitGroup
				wg.Add(20)
				postIds := make([]uint, 20)
				for i := 0; i < 20; i++ {
					go func(i int) {
						defer wg.Done()
						pid, err := postRepo.CreateNewPost(PostInfo{
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
					pE, err := postRepo.Storage.Read(postIds[i])
					Expect(err).Should(BeNil())
					Expect(pE.Content).Should(Equal("TestConcurrency" + strconv.Itoa(i)))
				}
			})
		})
	})
	Describe("Read Tweet", func() {
		Context("Read the same tweet 20 times concurrently", func() {
			It("should return 20 different pointers with the same information", func() {
				var wg sync.WaitGroup
				wg.Add(20)
				postIds := make([]uint, 20)
				for i := 0; i < 20; i++ {
					go func(i int) {
						defer wg.Done()
						defer GinkgoRecover()
						pid, err := postRepo.CreateNewPost(PostInfo{
							UserID:  puId,
							Content: "TestConcurrency",
						})
						Expect(err).To(BeNil())
						Expect(pid).To(Not(BeZero()))
						postIds[i] = pid
					}(i)
				}
				wg.Wait()
				postEs := make([]*storage.PostEntity, 20)
				wg.Add(20)
				for i := 0; i < 20; i++ {
					go func(i int) {
						defer wg.Done()
						post := postRepo.SelectById(postIds[i])
						postEs[i] = post
					}(i)
				}
				wg.Wait()
				prev := postRepo.SelectById(postIds[0])
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
