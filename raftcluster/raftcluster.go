package raftcluster

import (
    "flag"
    "github.com/coreos/etcd/raft/raftpb"
)

func StartRaftCluster() {
    //cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
    //id := flag.Int("id", 1, "node ID")
    //kvport := flag.Int("port", 9121, "key-value server port")
    //join := flag.Bool("join", false, "join an existing cluster")
    flag.Parse()

    proposeC := make(chan string)
    defer close(proposeC)
    confChangeC := make(chan raftpb.ConfChange)
    defer close(confChangeC)

    //var clusters [3]string
    //clusters := []string{CLUSTER_1_ADDR, CLUSTER_2_ADDR, CLUSTER_3_ADDR}
    // raft provides a commit stream for the proposals from the http api

    //getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }
    //getSnapshot := func() () {}
    //commitC, errorC, snapshotterReady := newRaftNode(*id, clusters, *join,
    //    getSnapshot, proposeC, confChangeC)

    //var kvs *kvstore
    //kvs = newKVStore(<-snapshotterReady, proposeC, commitC, errorC)
    //
    //// the key-value http handler will propose updates to raft
    //serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)
}
