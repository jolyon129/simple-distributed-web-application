package store

import (
    "bytes"
    "context"
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
    // Register types that will be transferred as implementations of interface values
    gob.Register(SessionParams{})
    gob.Register(SessionProviderParams{})
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

func (s *DBStore) GetSnapshot() ([]byte, error) {
    return s.MarshalJSON()
}

func (s *DBStore) MarshalJSON() ([]byte, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return json.Marshal(s)
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
    Key   string
    Value string
}

func NewStore(snapshotter *snap.Snapshotter, proposeC chan<- string,
        commitC <-chan *string, errorC <-chan error) *DBStore {
    s := &DBStore{
        proposeC:         proposeC,
        snapshotter:      snapshotter,
        commandIdCounter: 1000,
    }
    s.readCommits(commitC, errorC)
    go s.readCommits(commitC, errorC)
    return s
}

// Generate a new command ID for subscribe
func (s *DBStore) getCommandID() uint64 {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.commandIdCounter++
    return s.commandIdCounter
}

// Issue a propose to the raft cluster.
// The propose with cmdLog will be sent to a proposeC.
// The same proposeC will be consumed by the raft cluster.
// This method return immediately without waiting for the raft cluster
func (s *DBStore) propose(target string, params interface{}, id uint64) {
    var buf bytes.Buffer
    command := CommandLog{
        TargetMethod: target,
        ID:           id,
        Params:       params,
    }
    if err := gob.NewEncoder(&buf).Encode(command); err != nil {
        log.Fatal(err)
    }
    s.proposeC <- buf.String()
}


// Request a propose and wait till its committed and executed.
// This method should be called by the HTTP API controllers.
func (s *DBStore) RequestPropose(ctx context.Context, target string, params interface{}) (interface{},
        error) {
    cmdId := s.getCommandID()
    ch := managerSingle.subscribe(cmdId)
    resultC := ch.resultC
    errC := ch.errC
    defer close(resultC)  // Remember to close the channel
    defer close(errC)    // Remember to close the channel
    defer delete(managerSingle.proposeListener, cmdId)
    s.propose(target, params, cmdId)
    select {
    case result := <-resultC: //Wait till this propose commit
        return result, nil
    case err := <-errC:
        return nil, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (s *DBStore) readCommits(commitC <-chan *string, errorC <-chan error) {
    // Once raft commit the log, now we can perform the operation in log
    for command := range commitC {
        if command == nil { // The first command should be nil sent by raft
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
        // If the command is sent by clients
        // execute the command to the DBStore
        s.execute(command)
    }
    if err, ok := <-errorC; ok {
        log.Fatal(err)
    }
}

func (s *DBStore) recoverFromSnapshot(snapshot []byte) error {
    var persistentStore persistent
    if err := json.Unmarshal(snapshot, &persistentStore); err != nil {
        return err
    }
    s.mu.Lock()
    s.persistent = persistentStore
    s.mu.Unlock()
    return nil
}

// Execute a command.
func (s *DBStore) execute(command *string) {
    var cmdLog CommandLog
    dec := gob.NewDecoder(bytes.NewBufferString(*command))
    if err := dec.Decode(&cmdLog); err != nil {
        log.Fatalf("raftexample: could not decode message (%v)", err)
    }
    // How to Execute cmdLog?
    // Naive way
    switch cmdLog.TargetMethod {
    case METHOD_SessionInit:
        params, _ := cmdLog.Params.(SessionProviderParams)
        sessIns, err := sessProvider.SessionInit(params.Sid)
        managerSingle.notify(cmdLog.ID, sessIns, err) // Notify the corresponding listener to process
    }
}
