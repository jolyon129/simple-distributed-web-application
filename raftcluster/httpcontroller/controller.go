package httpcontroller

import (
    "github.com/coreos/etcd/raft/raftpb"
    "time"
    "zl2501-final-project/raftcluster/store"
)

const ContextTimeoutDuration = 5 * time.Second

var raftStore *store.DBStore
var confChangeC chan raftpb.ConfChange
var errC <-chan error

func InitController(db *store.DBStore,changeC chan raftpb.ConfChange, errC <-chan error) {
    raftStore = db
    confChangeC= changeC
    errC = errC
}



