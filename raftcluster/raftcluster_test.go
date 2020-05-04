package raftcluster_test

import (
    "context"
    "github.com/coreos/etcd/raft/raftpb"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "math/rand"
    "strconv"
    "strings"
    "sync"
    "zl2501-final-project/raftcluster/raft"
    "zl2501-final-project/raftcluster/store"
    authstorage "zl2501-final-project/raftcluster/store/authstore"
)

var _ = Describe("Raftcluster", func() {
    var DBStore *store.DBStore
    var proposeC chan string
    var confChangeC chan raftpb.ConfChange
    BeforeSuite(func() {
        proposeC = make(chan string)
        confChangeC = make(chan raftpb.ConfChange)
        getSnapshot := func() ([]byte, error) { return DBStore.GetSnapshot() }
        id := 1
        cluster := "http://127.0.0.1:9021"
        commitC, errorC, snapshotterReady := raft.NewRaftNode(id, strings.Split(cluster, ","),
            false, getSnapshot, proposeC, confChangeC)

        DBStore = store.NewStore(<-snapshotterReady, proposeC, commitC, errorC)
    })
    //AfterSuite(func() {
    //    close(proposeC)
    //    close(confChangeC)
    //})
    PDescribe("Request a propose to raft cluster", func() {
        Context("with session init request", func() {
            It("should return a new  session", func() {
                params := store.SessionProviderParams{
                    Sid: "fakesess12312" + string(rand.Intn(100000000)),
                }
                result, err := DBStore.RequestPropose(context.TODO(), store.METHOD_SessionInit,
                    params)
                sess, _ := result.(authstorage.SessionStorageInterface)
                Expect(err).Should(BeNil())
                Expect(sess.SessionID()).Should(Equal(params.Sid))
            })
        })
    })
    Describe("Test Concurrency", func() {
        Context("when multiple requests coming concurrently", func() {
            It("should be fine", func() {
                var wg sync.WaitGroup
                wg.Add(5)
                for i := 0; i < 5; i++ {
                    go func(i int) {
                        defer GinkgoRecover()
                        defer wg.Done()
                        params := store.SessionProviderParams{
                            Sid: "fakesess12312" + strconv.Itoa(rand.Intn(100000)+i),
                        }
                        result, err := DBStore.RequestPropose(context.Background(), store.METHOD_SessionInit,
                            params)
                        sess, _ := result.(authstorage.SessionStorageInterface)
                        //print(sess.SessionID())
                        Expect(err).Should(BeNil())
                        Expect(sess.SessionID()).Should(Equal(params.Sid))
                    }(i)
                }
                wg.Wait()
            })
        })
    })
})
