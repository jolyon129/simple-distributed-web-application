package store

import (
    "encoding/gob"
    "github.com/coreos/etcd/snap"
    "sync"
    "time"
    authmemory "zl2501-final-project/raftcluster/store/authstore/memory"
    beStorage "zl2501-final-project/raftcluster/store/backendstore"
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
    SrcId    uint
    TargetId uint
}
type UserStartFollowingDBParams struct {
    SrcId    uint
    TargetId uint
}
type UserStopFollowingDBParams struct {
    SrcId    uint
    TargetId uint
}

type TweetReadParams struct {
    TId uint
}
type TweetDeleteParams struct {
    TId uint
}
type TweetDeleteByCreatedTimeParams struct{
    TimeStamp time.Time
}

type UserAddTweetToUserParams struct {
    UId uint
    TId uint
}

// A memory store backed by raft cluster
type DBStore struct {
    mu               sync.RWMutex
    proposeC         chan<- string
    commandIdCounter uint64
    BeManager beStorage.Manager
    SessProvider authmemory.Provider
    snapshotter *snap.Snapshotter
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
    Sid string
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