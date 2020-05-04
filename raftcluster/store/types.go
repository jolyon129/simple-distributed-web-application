package store

import (
    "encoding/gob"
    "github.com/coreos/etcd/snap"
    "sync"
    "time"
    authmemory "zl2501-final-project/raftcluster/store/authstore/memory"
    bgmemory "zl2501-final-project/raftcluster/store/backendstore/memory"
)

const MaxLifeTime = 7200

func init() {
    // Register types that will be transferred as implementations of interface values
    // These struct wll be used when deserialize the commandLog
    gob.Register(SessionParams{})
    gob.Register(SessionProviderParams{})
    gob.Register(UserCreateParams{})
    gob.Register(UserIDParams{})
    gob.Register(UserInfo{})
    gob.Register(UserAddTweetToUserParams{})
    gob.Register(UserCheckWhetherFollowingDBParams{})
    gob.Register(UserStartFollowingDBParams{})
    gob.Register(UserStopFollowingDBParams{})
    gob.Register(TweetReadParams{})
    gob.Register(TweetDeleteParams{})
    gob.Register(TweetDeleteByCreatedTimeParams{})
    gob.Register(TweetInfo{})
}

type UserCheckWhetherFollowingDBParams struct {
    srcId    uint
    targetId uint
}
type UserStartFollowingDBParams struct {
    srcId    uint
    targetId uint
}
type UserStopFollowingDBParams struct {
    srcId    uint
    targetId uint
}

type TweetReadParams struct {
    tId uint
}
type TweetDeleteParams struct {
    tId uint
}
type TweetDeleteByCreatedTimeParams struct{
    timeStamp time.Time
}

type UserAddTweetToUserParams struct {
    uId uint
    tId uint
}

type DBStore struct {
    mu               sync.RWMutex
    proposeC         chan<- string
    commandIdCounter uint64
    persistent
    snapshotter *snap.Snapshotter
}
type persistent struct {
    TweetStore   bgmemory.MemTweetStore
    UserStore    bgmemory.MemUserStore
    SessionStore authmemory.MemSessStore
}

// CommandLog for serialize and deserialize
type CommandLog struct {
    TargetMethod string
    ID           uint64
    Params       interface{}
}

type SessionProviderParams struct {
    Sid string
}

type SessionParams struct {
    SessionProviderParams
    Key   string
    Value string
}

type UserCreateParams struct {
    UserInfo
}

type UserIDParams struct {
    Uid uint
}

type UserInfo struct {
    UserName string
    Password string
}


type TweetInfo struct {
    UserID  uint
    Content string
}