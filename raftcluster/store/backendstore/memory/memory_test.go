package memory_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "zl2501-final-project/raftcluster/store/backendstore"
    _ "zl2501-final-project/raftcluster/store/backendstore/memory"
)

var _ = Describe("Memory", func() {
    var storageManager *Manager
    BeforeSuite(func() {
        storageManager = NewManager("memory")
    })
    Describe("Get the snapshot of UserStorage", func() {
        Context("when have some users info inside", func() {
            It("should return the right json", func() {
                resultC := make(chan uint, 10)
                errorC := make(chan error, 10)
                storageManager.UserStorage.Create(&UserEntity{
                    UserName: "jolyon129",
                    Password: "123",
                }, resultC, errorC)
                userId1 := <-resultC
                storageManager.UserStorage.Create(&UserEntity{
                    UserName: "zhuolun",
                    Password: "123",
                }, resultC, errorC)
                userId2 := <-resultC
                resultC1 := make(chan bool, 10)
                storageManager.UserStorage.AddTweetToUserDB(userId1, 20003, resultC1, errorC)
                storageManager.UserStorage.AddTweetToUserDB(userId1, 20004, resultC1, errorC)
                storageManager.UserStorage.StartFollowingDB(userId1, userId2, resultC1, errorC)
                storageManager.UserStorage.StartFollowingDB(userId2, userId1, resultC1, errorC)
                j, err := storageManager.UserStorage.GetSnapshot()
                Expect(err).Should(BeNil())
                js := string(j)
                //println(js)
                Expect(js).ShouldNot(BeEmpty())
                Expect(j).Should(MatchJSON([]byte(
                        `{"pkCounter":102,"userMap":{"101":{"Follower":[102],"Following":[102],"ID":101,"Password":"123","Tweets":[20003,20004],"UserName":"jolyon129"},"102":{"Follower":[101],"Following":[101],"ID":102,"Password":"123","Tweets":[],"UserName":"zhuolun"}},"userNameSet":{"jolyon129":true,"zhuolun":true}}`)))
            })
        })
    })
    Describe("Get the snapshot of TweetStorage", func() {
        Context("when have some users info inside", func() {
            It("should return the right json", func() {
                resultC := make(chan uint, 10)
                errorC := make(chan error, 10)
                storageManager.TweetStorage.Create(&TweetEntity{
                    UserID:      1001,
                    Content:     "Test",
                }, resultC, errorC)
                storageManager.TweetStorage.Create(&TweetEntity{
                    UserID:      1001,
                    Content:     "Test2",
                }, resultC, errorC)
                j, err := storageManager.TweetStorage.GetSnapshot()
                Expect(err).Should(BeNil())
                js := string(j)
                //println(js)
                Expect(js).Should(ContainSubstring(`"Content":"Test"`))
                Expect(js).Should(ContainSubstring(`"Content":"Test2"`))
            })
        })
    })
})
