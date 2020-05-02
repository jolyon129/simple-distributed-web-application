package store

import (
    "encoding/json"
    "sync"
    authmemory "zl2501-final-project/raftcluster/store/authstore/memory"
    bgmemory "zl2501-final-project/raftcluster/store/backendstore/memory"
)

type store struct {
    mu sync.RWMutex
    persistent
}

type persistent struct {
    TweetStore   bgmemory.MemTweetStore
    UserStore    bgmemory.MemUserStore
    SessionStore authmemory.MemSessStore
}

func (s *store) getSnapshot() ([]byte, error) {
    return s.MarshalJSON()
}

func (s *store) MarshalJSON() ([]byte, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return json.Marshal(s.persistent)
}
