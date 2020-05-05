package raftcluster_test

import (
    "context"
    "encoding/json"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "math/rand"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "sync"
    "zl2501-final-project/raftcluster/store"
    authstorage "zl2501-final-project/raftcluster/store/authstore"
)

var _ = Describe("Raftcluster", func() {
    PDescribe("Request a propose to raft cluster", func() {
        var DBStore *store.DBStore
        //var proposeC chan string
        //var confChangeC chan raftpb.ConfChange
        //BeforeSuite(func() {
        //    cluster := flag.String("cluster", "http://127.0.0.1:9011", "comma separated cluster peers")
        //    id := flag.Int("id", 1, "node ID")
        //    //httpAPIPort := flag.Int("port", 9001, "key-value server port")
        //    join := flag.Bool("join", false, "join an existing cluster")
        //    flag.Parse()
        //
        //    proposeC = make(chan string)
        //    defer close(proposeC)
        //    confChangeC = make(chan raftpb.ConfChange)
        //    defer close(confChangeC)
        //    getSnapshot := func() ([]byte, error) { return DBStore.GetSnapshot() }
        //
        //    commitC, errorC, snapshotterReady := raft.NewRaftNode(*id, strings.Split(*cluster, ","),
        //        *join, getSnapshot, proposeC, confChangeC)
        //
        //    DBStore = store.NewStore(<-snapshotterReady, proposeC, commitC, errorC)
        //})
        //AfterSuite(func() {
        //    if proposeC!=nil{
        //        close(proposeC)
        //    }
        //    if confChangeC!=nil{
        //        close(confChangeC)
        //    }
        //})
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
    Describe("Test http API", func() {
        PContext("session init", func() {
            It("return sessid and success", func() {
                sid := "fakesessid12321" + strconv.Itoa(rand.Int())
                resp, err := http.PostForm("http://127.0.0.1:9001/session/"+sid, nil)
                Expect(err).Should(BeNil())
                defer resp.Body.Close()
                var result map[string]interface{}
                json.NewDecoder(resp.Body).Decode(&result)
                Expect(result["result"]).Should(Equal(sid))
            })
        })
        PContext("change request URL", func() {
            It("change url", func() {
                sid := "fakesessid12321" + strconv.Itoa(rand.Int())
                r, _ := http.NewRequest("POST", "/session/"+sid, nil)
                newReq, _ := http.NewRequest(r.Method, "http://127.0.0.1:9001"+r.URL.Path, r.Body)
                resp, err := http.DefaultClient.Do(newReq)
                //resp, err := http.PostForm("http://127.0.0.1:9001/session/"+sid, nil)
                Expect(err).Should(BeNil())
                defer resp.Body.Close()
                var result map[string]interface{}
                json.NewDecoder(resp.Body).Decode(&result)
                Expect(result["result"]).Should(Equal(sid))
            })
        })
        Context("custom put request with data", func() {
            It("should be fine", func() {
                sid := "fakesessid12321" + strconv.Itoa(rand.Int())
                data := url.Values{}
                data.Set("value", "jolyon129")
                r, _ := http.NewRequest("PUT", "/session/"+sid+"/name",
                    strings.NewReader(data.Encode()))
                newReq, _ := http.NewRequest(r.Method, "http://127.0.0.1:9001"+r.URL.Path, r.Body)
                newReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
                resp, err := http.DefaultClient.Do(newReq)
                //resp, err := http.PostForm("http://127.0.0.1:9001/session/"+sid, nil)
                Expect(err).Should(BeNil())
                defer resp.Body.Close()
                var result map[string]interface{}
                json.NewDecoder(resp.Body).Decode(&result)
                Expect(result["result"]).Should(BeTrue())
            })
        })
    })
})
