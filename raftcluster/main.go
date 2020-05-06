package raftcluster

import (
    "flag"
    "github.com/coreos/etcd/raft/raftpb"
    "strings"
    "zl2501-final-project/raftcluster/raft"
    "zl2501-final-project/raftcluster/store"
)

func StartRaftCluster() {
    cluster := flag.String("cluster", "http://127.0.0.1:9011", "comma separated cluster peers")
    id := flag.Int("id", 1, "node ID")
    httpAPIPort := flag.Int("port", 9004, "key-value server port")
    join := flag.Bool("join", false, "join an existing cluster")
    flag.Parse()

    var DBStore *store.DBStore
    proposeC := make(chan string)
    defer close(proposeC)
    confChangeC := make(chan raftpb.ConfChange)
    defer close(confChangeC)
    getSnapshot := func() ([]byte, error) { return DBStore.GetSnapshot() }

    commitC, errorC, snapshotterReady := raft.NewRaftNode(*id, strings.Split(*cluster, ","),
        *join, getSnapshot, proposeC, confChangeC)

    DBStore = store.NewStore(<-snapshotterReady, proposeC, commitC, errorC)
    ServerHttpAPI(DBStore, *httpAPIPort, confChangeC, errorC)

}
