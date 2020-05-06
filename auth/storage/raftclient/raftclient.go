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
    "zl2501-final-project/auth/storage"
)

const (
    ADDR1 = "http://127.0.0.1:9004"
    ADDR2 = "http://127.0.0.1:9005"
    ADDR3 = "http://127.0.0.1:9006"
)

func init() {
    storage.RegisterProvider("raft", pder)
}

var pder = &RaftSessProvider{
    addrs: []string{
        ADDR1, ADDR2, ADDR3,
    },
    client: &http.Client{
        Timeout: 2 * time.Second,
    },
}

// This is an abstract provider storage
type RaftSessProvider struct {
    addrs  []string
    client *http.Client
}

// This is an abstract session storage
type RaftSessStore struct {
    sid string
}

type SessionValueSetResultType struct {
    Error  string
    Result bool
}

func (r RaftSessStore) Set(key, value interface{}) error {
    form := url.Values{}
    switch value.(type) {
    case int:
        form.Set("value", strconv.Itoa(value.(int)))
    case string:
        form.Set("value", value.(string))
    }
    req, _ := http.NewRequest("PUT", "/session/"+r.sid+"/"+key.(string),
        strings.NewReader(form.Encode()))
    resp, err := doRequestWithRetry(req)
    defer resp.Body.Close()
    if err != nil {
        return err
    }
    //body, _ := ioutil.ReadAll(resp.Body)
    //var result SessionValueSetResultType
    //json.Unmarshal(body, &result)
    return nil
}

type SessionValueGetResultType struct {
    Error  string
    Result string
}

func (r RaftSessStore) Get(key interface{}) interface{} {
    form := url.Values{}
    req, _ := http.NewRequest("GET", "/session/"+r.sid+"/"+key.(string),
        strings.NewReader(form.Encode()))
    resp, err := doRequestWithRetry(req)
    defer resp.Body.Close()
    if err != nil {
        return err
    }
    var result SessionValueGetResultType
    body, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(body, &result)
    if result.Error != "" {
        return nil
    }
    return result.Result
}

func (r RaftSessStore) Delete(key interface{}) error {
    panic("implement me")
}

func (r RaftSessStore) SessionID() string {
    return r.sid
}

// This wrapper will try all raft nodes one by one till one of them
// responses.
func doRequestWithRetry(r *http.Request) (*http.Response,
        error) {
    var resp *http.Response
    var err error
    var index int
    for idx, addr := range pder.addrs {
        // Create a new request with raft node host address
        newReq, _ := http.NewRequest(r.Method, addr+r.URL.Path, r.Body)
        if r.Method == "PUT" || r.Method == "POST" {
            newReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        }
        resp, err = pder.client.Do(newReq)
        if err != nil {
            continue
        }
        index = idx
        break
    }
    // Move the current working raft node address to the first one
    if index != 0 {
        old := pder.addrs[0]
        addr := pder.addrs[index]
        pder.addrs[0] = addr
        pder.addrs[index] = old
    }
    return resp, err
}

func (p RaftSessProvider) SessionInit(sid string) (storage.SessionStorageInterface, error) {
    req, _ := http.NewRequest("POST", "/session/"+sid, nil)
    resp, err := doRequestWithRetry(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    //println(body)
    var result sessionSidResult
    json.Unmarshal(body, &result)
    if result.Error != "" {
        return nil, errors.New(result.Error)
    }
    sess := &RaftSessStore{sid: result.Result}
    return sess, err
}

type sessionSidResult struct {
    Result string `json:"result"`
    Error  string `json:"error"`
}

type MemSessStore struct {
    Sid          string                      // unique session id
    TimeAccessed time.Time                   // last access time
    Value        map[interface{}]interface{} // session value stored inside
}

type sessionReadResult struct {
    Error  string
    Result MemSessStore
}

func (p RaftSessProvider) SessionRead(sid string) (storage.SessionStorageInterface, error) {
    req, _ := http.NewRequest("GET", "/session/"+sid, nil)
    resp, err := doRequestWithRetry(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    //ioutil.ReadAll(res.Body)
    var result sessionReadResult
    //json.NewDecoder(resp.Body).Decode(&result)
    body, _ := ioutil.ReadAll(resp.Body)
    //println(body)
    //tesJ := `{"error":"the session Id: 1j28v6loBj65ypDacf5VJxRDXDcDRU8y1RkdXNOu4qo=5577006791947779410 is not existed","result":null}`
    json.Unmarshal(body, &result)
    if result.Error != "" {
        return nil, errors.New(result.Error)
    }
    sess := &RaftSessStore{sid: result.Result.Sid}
    return sess, err
}

func (p RaftSessProvider) SessionDestroy(sid string) error {
    req, _ := http.NewRequest("DELETE", "/session/"+sid, nil)
    resp, err := doRequestWithRetry(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func (p RaftSessProvider) SessionGC(maxLifeTime int64) {
    req, _ := http.NewRequest("POST", "/sessiongc", nil)
    doRequestWithRetry(req)
}
