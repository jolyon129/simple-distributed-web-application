package store

import (
    "bytes"
    "context"
    "encoding/gob"
    "encoding/json"
    "fmt"
    "github.com/coreos/etcd/snap"
    "log"
    "strings"
    authstore "zl2501-final-project/raftcluster/store/authstore"
    _ "zl2501-final-project/raftcluster/store/authstore/memory"
    "zl2501-final-project/raftcluster/store/backendstore"
    _ "zl2501-final-project/raftcluster/store/backendstore/memory"
)

var sessProvider authstore.ProviderInterface
var bkStorageManager *backendstore.Manager

func init() {
    sessProvider, _ = authstore.GetProvider("memory")
    bkStorageManager = backendstore.NewManager("memory")
}

//const ProviderName = "memory"
//const CookieName = "newSessionId"

func (s *DBStore) GetSnapshot() ([]byte, error) {
    return s.MarshalJSON()
}

func (s *DBStore) MarshalJSON() ([]byte, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return json.Marshal(s)
}

func NewStore(snapshotter *snap.Snapshotter, proposeC chan<- string,
        commitC <-chan *string, errorC <-chan error) *DBStore {
    s := &DBStore{
        proposeC:         proposeC,
        snapshotter:      snapshotter,
        commandIdCounter: 1000,
    }
    s.readCommits(commitC, errorC) // replay log into DBStore
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
func (s *DBStore) RequestPropose(ctx context.Context, targetMethod string, params interface{}) (interface{},
        error) {
    cmdId := s.getCommandID()
    ch := managerSingle.subscribe(cmdId)
    defer managerSingle.unsubscribe(cmdId)
    resultC := ch.resultC
    errC := ch.errC
    defer close(resultC) // Remember to close the channel
    defer close(errC)    // Remember to close the channel

    go func() {
        // If its a get function,
        // no need to propose to raft
        // execute immediately
        if strings.Contains(targetMethod, "Get") || strings.Contains(targetMethod, "Read") {
            log.Printf("%s is a GET command, no need to propose. Execute immediately", targetMethod)
            s.execute(CommandLog{
                TargetMethod: targetMethod,
                ID:           cmdId,
                Params:       params,
            })
        } else {
            s.propose(targetMethod, params, cmdId)
        }
    }()

    select {
    case result := <-resultC: //Wait till this propose executed
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
        // execute the command to the DBStore
        cmdIns := deserialize(command)
        s.execute(*cmdIns)
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

func deserialize(command *string) *CommandLog {
    var cmdLog CommandLog
    dec := gob.NewDecoder(bytes.NewBufferString(*command))
    if err := dec.Decode(&cmdLog); err != nil {
        log.Fatalf("raftexample: could not decode message (%v)", err)
    }
    return &cmdLog
}

// Execute a command.
func (s *DBStore) execute(cmdLog CommandLog) {
    // How to Execute cmdLog?
    // Naive way
    switch cmdLog.TargetMethod {
    case METHOD_SessionInit:
        params, _ := cmdLog.Params.(SessionProviderParams)
        sessIns, err := sessProvider.SessionInit(params.Sid)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
        } else {
            managerSingle.notify(cmdLog.ID, sessIns)
        }
    case METHOD_SessionRead:
        params, _ := cmdLog.Params.(SessionProviderParams)
        sessIns, err := sessProvider.SessionRead(params.Sid)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
        } else {
            managerSingle.notify(cmdLog.ID, sessIns)
        }
    case METHOD_SessionDestroy:
        params, _ := cmdLog.Params.(SessionProviderParams)
        err := sessProvider.SessionDestroy(params.Sid)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
        } else {
            managerSingle.notify(cmdLog.ID, true)
        }
    case METHOD_SessionGC:
        sessProvider.SessionGC(MaxLifeTime)
        managerSingle.notify(cmdLog.ID, true)
    case METHOD_SessionGet:
        params, _ := cmdLog.Params.(SessionParams)
        sessIns, err := sessProvider.SessionRead(params.Sid)
        value := sessIns.Get(params.Key)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
        } else {
            managerSingle.notify(cmdLog.ID, value)
        }
    case METHOD_SessionSet:
        params, _ := cmdLog.Params.(SessionParams)
        sessIns, err := sessProvider.SessionRead(params.Sid)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        err1 := sessIns.Set(params.Key, params.Value)
        if err1 != nil {
            managerSingle.notifyError(cmdLog.ID, err1)
            return
        }
        managerSingle.notify(cmdLog.ID, true)
    case METHOD_SessionDelete:
        params, _ := cmdLog.Params.(SessionParams)
        sessIns, err := sessProvider.SessionRead(params.Sid)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        err1 := sessIns.Delete(params.Key)
        if err1 != nil {
            managerSingle.notifyError(cmdLog.ID, err1)
            return
        }
        managerSingle.notify(cmdLog.ID, true)
    case METHOD_UserCreate:
        params, _ := cmdLog.Params.(UserCreateParams)
        uID, err := CreateNewUser(context.Background(), &params.UserInfo)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, uID)
    case METHOD_UserDelete:
        managerSingle.notifyError(cmdLog.ID,
            fmt.Errorf("this method %s has not been implemented", cmdLog.TargetMethod))
    case METHOD_UserGet:
        params, _ := cmdLog.Params.(UserIDParams)
        userEnt, err := UserSelectById(context.Background(), params.Uid)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, userEnt)
    case METHOD_UserUpdate:
        fmt.Errorf("this method %s has not been implemented", cmdLog.TargetMethod)
    case METHOD_UserFindAll:
        userEnts, err := FindAllUsers(context.Background())
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, userEnts)
    case METHOD_UserAddTweetToUserDB:
        params, _ := cmdLog.Params.(UserAddTweetToUserParams)
        ok, err := AddTweetToUser(context.Background(), params.UId, params.TId)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, ok)
    case METHOD_UserCheckWhetherFollowingGetDB:
        params, _ := cmdLog.Params.(UserCheckWhetherFollowingDBParams)
        ret, err := CheckWhetherFollowing(context.Background(), params.SrcId, params.TargetId)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, ret)
    case METHOD_UserStartFollowingDB:
        params, _ := cmdLog.Params.(UserStartFollowingDBParams)
        ok, err := StartFollowing(context.Background(), params.SrcId, params.TargetId)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, ok)
    case METHOD_UserStopFollowingDB:
        params, _ := cmdLog.Params.(UserStopFollowingDBParams)
        ok, err := StopFollowing(context.Background(), params.SrcId, params.TargetId)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, ok)
    case METHOD_TweetCreate:
        params, _ := cmdLog.Params.(TweetInfo)
        tId, err := SaveTweet(context.Background(), params)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, tId)
    case METHOD_TweetGet:
        params, _ := cmdLog.Params.(TweetReadParams)
        tweet, err := TweetSelectById(context.Background(), params.TId)
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, tweet)
    case METHOD_TweetGetAll:
        tweets, err := TweetGetAllTweets(context.Background())
        if err != nil {
            managerSingle.notifyError(cmdLog.ID, err)
            return
        }
        managerSingle.notify(cmdLog.ID, tweets)
    case METHOD_TweetDelete:
        fmt.Errorf("this method %s has not been implemented", cmdLog.TargetMethod)
        //params, _ := cmdLog.Params.(TweetDeleteParams)
        //ok, err := DeleteById(context.Background(), params.tId)
        //if err != nil {
        //    managerSingle.notifyError(cmdLog.ID, err)
        //    return
        //}
        //managerSingle.notify(cmdLog.ID, ok)
    case METHOD_TweetDeleteByCreatedTime:
        fmt.Errorf("this method %s has not been implemented", cmdLog.TargetMethod)
        //params, _ := cmdLog.Params.(TweetDeleteByCreatedTimeParams)
        //ok, err := TweetDeleteByCreatedTime(context.Background(), params.timeStamp)
        //if err != nil {
        //    managerSingle.notifyError(cmdLog.ID, err)
        //    return
        //}
        //managerSingle.notify(cmdLog.ID, ok)
    }
}
