package httpcontroller

import (
    "github.com/coreos/etcd/raft/raftpb"
    "zl2501-final-project/raftcluster/store"
)


var raftStore *store.DBStore
var confChangeC chan raftpb.ConfChange
var errC <-chan error

func InitController(db *store.DBStore,changeC chan raftpb.ConfChange, errC <-chan error) {
    raftStore = db
    confChangeC= changeC
    errC = errC
}



