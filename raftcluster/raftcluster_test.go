package raftcluster_test

import (
    "context"
    "github.com/coreos/etcd/raft/raftpb"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "math/rand"
    "strings"
    "zl2501-final-project/raftcluster/raft"
    "zl2501-final-project/raftcluster/store"
    authstorage "zl2501-final-project/raftcluster/store/authstore"
)

var _ = Describe("Raftcluster", func() {
    var MemStore *store.DBStore
    var proposeC chan string
    var confChangeC chan raftpb.ConfChange
    BeforeSuite(func() {
        proposeC = make(chan string)
        confChangeC = make(chan raftpb.ConfChange)
        getSnapshot := func() ([]byte, error) { return MemStore.GetSnapshot() }
        id := 1
        cluster := "http://127.0.0.1:9021"
        commitC, errorC, snapshotterReady := raft.NewRaftNode(id, strings.Split(cluster, ","),
            false, getSnapshot, proposeC, confChangeC)

        MemStore = store.NewStore(<-snapshotterReady, proposeC, commitC, errorC)
    })
    AfterSuite(func() {
        close(proposeC)
        close(confChangeC)
    })
    Describe("Request a propose to raft cluster", func() {
        Context("with session init request", func() {
            It("should return a new  session", func() {
                params := store.SessionProviderParams{
                    Sid: "fakesess12312" + string(rand.Int()),
                }
                result, err := MemStore.RequestPropose(context.TODO(), store.METHOD_SessionInit,
                    params)
                sess, _ := result.(authstorage.SessionStorageInterface)
                Expect(err).Should(BeNil())
                Expect(sess.SessionID()).Should(Equal(params.Sid))
                //mStore.RequestPropose(context.TODO(), METHOD_SessionRead, params1)
            })
        })
    })
})
