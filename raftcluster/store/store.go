package store

import (
    "bytes"
    "encoding/gob"
    "encoding/json"
    "github.com/coreos/etcd/snap"
    "log"
    "sync"
    authstorage "zl2501-final-project/raftcluster/store/authstore"
    authmemory "zl2501-final-project/raftcluster/store/authstore/memory"
    bgmemory "zl2501-final-project/raftcluster/store/backendstore/memory"
)

var sessProvider authstorage.ProviderInterface

func init() {
    sessProvider, _ = authstorage.GetProvider("memory")
}

type store struct {
    mu       sync.RWMutex
    proposeC chan<- string
    persistent
    snapshotter *snap.Snapshotter
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
    return json.Marshal(s)
}

type commandLog struct {
    target string
    method string
    ID     uint64
    params interface{}
}

type SessionProviderParams struct {
    sid string
}

type SessionParams struct {
    key   string
    value string
}

func NewStore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *string, errorC <-chan error) {
    s := &store{
        proposeC:    proposeC,
        snapshotter: snapshotter,
    }
    s.readCommits(commitC, errorC)
}

// Issue a propose with command and send to proposeC
// The same proposeC will be consumed by the raft
func (s *store) propose(command commandLog) {
    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(command); err != nil {
        log.Fatal(err)
    }
    s.proposeC <- buf.String()
}

func (s *store) readCommits(commitC <-chan *string, errorC <-chan error) {
    // Once raft commit the log, now we can perform the operation in log
    for command := range commitC {
        if command == nil {
            // done replaying log; new command incoming
            // OR signaled to load snapshot
            snapshot, err := s.snapshotter.Load()
            if err == snap.ErrNoSnapshot {
                return
            }
            if err != nil && err != snap.ErrNoSnapshot {
                log.Panic(err)
            }
            log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
            if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
                log.Panic(err)
            }
            continue
        }

        s.notifyProposeEventManager()

        //
        //var dataKv kv
        //dec := gob.NewDecoder(bytes.NewBufferString(*command))
        //if err := dec.Decode(&dataKv); err != nil {
        //    log.Fatalf("raftexample: could not decode message (%v)", err)
        //}
        //s.mu.Lock()
        //s.kvStore[dataKv.Key] = dataKv.Val
        //s.mu.Unlock()
    }
    if err, ok := <-errorC; ok {
        log.Fatal(err)
    }
}

func (s *store) recoverFromSnapshot(snapshot []byte) error {
    var persistentStore persistent
    if err := json.Unmarshal(snapshot, &persistentStore); err != nil {
        return err
    }
    s.mu.Lock()
    s.persistent = persistentStore
    s.mu.Unlock()
    return nil
}

func (s *store) notifyProposeEventManager()  {

}