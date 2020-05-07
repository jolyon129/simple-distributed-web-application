package raftclient

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
    "zl2501-final-project/backend/model/storage"
)

const (
    ADDR1 = "http://127.0.0.1:9004"
    ADDR2 = "http://127.0.0.1:9005"
    ADDR3 = "http://127.0.0.1:9006"
)

var Addrs []string
var Client *http.Client

func init() {
    var raftUserStorage storage.UserStorageInterface
    var raftTweetStorage storage.TweetStorageInterface
    Addrs = []string{
        ADDR1, ADDR2, ADDR3,
    }
    Client = &http.Client{
        Timeout: 2 * time.Second,
    }
    raftUserStorage = RaftUserStorage{}
    raftTweetStorage = RaftTweetStorage{}
    memModels := storage.Manager{
        UserStorage:  raftUserStorage,
        TweetStorage: raftTweetStorage,
    }
    storage.RegisterDriver("raft", &memModels)
}

type RaftUserStorage struct {
}
type RaftTweetStorage struct {
}

type CreateRetType struct {
    Error  string
    Result uint
}

func (r RaftTweetStorage) Create(tweet *storage.TweetEntity, resultC chan uint,
        errorChan chan error) uint {
    form := url.Values{}
    form.Set("uid", strconv.Itoa(int(tweet.UserID)))
    form.Set("content", tweet.Content)
    req, _ := http.NewRequest("POST", "/tweet/", strings.NewReader(form.Encode()))
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return 0
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
        return 0
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result CreateRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
        return 0
    }
    resultC <- result.Result
    return result.Result
}

type TweetReadRetType struct {
    Error  string
    Result storage.TweetEntity
}

func (r RaftTweetStorage) Read(ID uint, resultC chan *storage.TweetEntity, errorChan chan error) {
    req, _ := http.NewRequest("GET", "/tweet/"+strconv.Itoa(int(ID)), nil)
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
    }

    body, _ := ioutil.ReadAll(resp.Body)
    var result TweetReadRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }
    resultC <- &result.Result
}

func (r RaftTweetStorage) Delete(ID uint, result chan bool, errorChan chan error) {
    panic("implement me")
}

func (r RaftTweetStorage) DeleteByCreatedTime(timeStamp time.Time, result chan bool, errorChan chan error) {
    panic("implement me")
}

func (r RaftUserStorage) Create(user *storage.UserEntity, resultC chan uint, errorChan chan error) {
    form := url.Values{}
    form.Set("username", user.UserName)
    form.Set("password", user.Password)
    req, _ := http.NewRequest("POST", "/user", strings.NewReader(form.Encode()))
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result CreateRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }
    resultC <- result.Result
}

func (r RaftUserStorage) Delete(ID uint, result chan bool, errorChan chan error) {
    panic("implement me")
}

type UserReadRetType struct {
    Error  string
    Result storage.UserEntity
}

func (r RaftUserStorage) Read(ID uint, resultC chan *storage.UserEntity, errorChan chan error) {
    req, _ := http.NewRequest("GET", "/user/"+strconv.Itoa(int(ID)), nil)
    resp, err := doRequestWithRetry(req)

    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result UserReadRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }

    resultC <- &result.Result
}

func (r RaftUserStorage) Update(ID uint, user *storage.UserEntity, resultC chan uint,
        errorChan chan error) {
    panic("implement me")
}

type UserFindAllRetType struct {
    Error string
    Users []*storage.UserEntity
}

func (r RaftUserStorage) FindAll(resultC chan []*storage.UserEntity, errorChan chan error) {
    req, _ := http.NewRequest("GET", "/user", nil)
    resp, err := doRequestWithRetry(req)
    if resp == nil { // If response is null
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()
    if err != nil {
        errorChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result UserFindAllRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }
    resultC <- result.Users
}

type BoolRetType struct {
    Error  string
    Result bool
}

func (r RaftUserStorage) AddTweetToUserDB(uId uint, pId uint, resultC chan bool,
        errorChan chan error) {
    req, _ := http.NewRequest("POST", "/user/"+strconv.Itoa(int(uId))+"/tweet/"+strconv.Itoa(int(
        pId)), nil)
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result BoolRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }
    resultC <- result.Result
}

func (r RaftUserStorage) CheckWhetherFollowingDB(srcId uint, targetId uint,
        resultC chan bool, errChan chan error) {
    req, _ := http.NewRequest("GET", "/user/"+strconv.Itoa(int(srcId))+"/following/"+strconv.Itoa(
        int(targetId)), nil)
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result BoolRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errChan <- errors.New(result.Error)
    }
    resultC <- result.Result
}

func (r RaftUserStorage) StartFollowingDB(srcId uint, targetID uint,
        resultC chan bool, errorChan chan error) {
    req, _ := http.NewRequest("POST", "/user/"+strconv.Itoa(int(srcId))+"/following/"+strconv.Itoa(
        int(
            targetID)), nil)
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result BoolRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }
    resultC <- result.Result
}

func (r RaftUserStorage) StopFollowingDB(srcId uint, targetID uint, resultC chan bool,
        errorChan chan error) {
    req, _ := http.NewRequest("DELETE", "/user/"+strconv.Itoa(int(srcId))+"/following/"+strconv.Itoa(
        int(targetID)), nil)
    resp, err := doRequestWithRetry(req)
    if resp == nil {
        if err != nil {
            errorChan <- err
        }
        return
    }
    defer resp.Body.Close()

    if err != nil {
        errorChan <- err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var result BoolRetType
    json.Unmarshal(body, &result)
    if result.Error != "" {
        errorChan <- errors.New(result.Error)
    }
    resultC <- result.Result
}

// This wrapper will try all raft nodes one by one till one of them
// responses.
func doRequestWithRetry(r *http.Request) (*http.Response,
        error) {
    var resp *http.Response
    var err error
    var index int
    for idx, addr := range Addrs {
        // Create a new request with raft node host address
        newReq, _ := http.NewRequest(r.Method, addr+r.URL.Path, r.Body)
        if r.Method == "PUT" || r.Method == "POST" {
            newReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        }
        resp, err = Client.Do(newReq)
        if err != nil {
            continue
        }
        index = idx
        break
    }
    // Move the current working raft node address to the first one
    if index != 0 {
        old := Addrs[0]
        addr := Addrs[index]
        Addrs[0] = addr
        Addrs[index] = old
    }
    return resp, err
}
