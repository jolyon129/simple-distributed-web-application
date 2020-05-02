package memory

import (
    "container/list"
    "encoding/json"
    "errors"
    "sort"
    "sync"
    "zl2501-final-project/raftcluster/store/backendstore"
)

type MemUserStore struct {
    backendstore.UserStorageInterface
    sync.Mutex
    userMap     map[uint]*backendstore.UserEntity // Map index to entity/record
    userNameSet map[string]bool                   // A set of username. Used as a fast approach to avoid duplicate names
    pkCounter   uint                              // Primary Key Counter
}

// Sort User By their ID
type EntityById []*backendstore.UserEntity

func (e EntityById) Len() int {
    return len(e)
}

func (e EntityById) Less(i, j int) bool {
    return e[i].ID < e[j].ID
}

func (e EntityById) Swap(i, j int) {
    e[i], e[j] = e[j], e[i]
}

func (m *MemUserStore) FindAll(result chan []*backendstore.UserEntity, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    res := make(EntityById, 0)
    for _, entity := range m.userMap {
        newE := backendstore.UserEntity{}
        copyUserEntity(&newE, entity)
        //newList.PushBack(&newE)
        res = append(res, &newE)
    }
    sort.Sort(res) // Sort entity by their ID
    result <- res
}

// Return a new primary key
// This function does not need to be locked
func (m *MemUserStore) getNewPK() uint {
    m.pkCounter++
    return m.pkCounter
}

// Update can only modified the password and the post, the following and the follower list.
// Take O(#post+#following+#follower) time.
func (m *MemUserStore) Update(ID uint, user *backendstore.UserEntity, result chan uint,
        errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    totalNewE := &backendstore.UserEntity{}
    copyUserEntity(totalNewE, user)
    m.userMap[ID] = totalNewE
    result <- ID
}

func (m *MemUserStore) Create(user *backendstore.UserEntity, result chan uint, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if _, ok := m.userNameSet[user.UserName]; ok {
        errorChan <- &backendstore.MyStorageError{Message: "Duplicate UserStorage Names!"}
        return // Return Immediately!
    }
    pk := m.getNewPK()
    newUser := backendstore.UserEntity{
        ID:        pk,
        UserName:  user.UserName,
        Password:  user.Password,
        Tweets:    list.New(),
        Following: list.New(),
        Follower:  list.New(),
    }
    m.userMap[pk] = &newUser
    m.userNameSet[user.UserName] = true
    result <- pk
    return
}

func (m *MemUserStore) Delete(ID uint, result chan bool, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if _, ok := m.userMap[ID]; !ok {
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist ID"}
        return
    } else {
        uInDB := m.userMap[ID]
        // Copy the post list
        newUser := backendstore.UserEntity{}
        copyUserEntity(&newUser, uInDB)
        delete(m.userNameSet, newUser.UserName)
        delete(m.userMap, ID)
        result <- true
    }
}

func (m *MemUserStore) Read(ID uint, result chan *backendstore.UserEntity, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if _, ok := m.userMap[ID]; !ok {
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist ID"}
        return
    } else {
        uInDB := m.userMap[ID]
        newUser := backendstore.UserEntity{}
        copyUserEntity(&newUser, uInDB)
        result <- &newUser
    }
}

func (m *MemUserStore) AddTweetToUserDB(uId uint, pId uint, result chan bool, errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if _, ok := m.userMap[uId]; !ok {
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist ID"}
        return
    }
    uInDB := m.userMap[uId]
    newUser := backendstore.UserEntity{}
    copyUserEntity(&newUser, uInDB)
    newUser.Tweets.PushBack(pId)
    m.userMap[uId] = &newUser
    result <- true
}

// Check whether srcId is following targetId
func (m *MemUserStore) CheckWhetherFollowingDB(srcId uint, targetId uint, result chan bool,
        errChan chan error) {
    m.Lock()
    defer m.Unlock()
    _, ok2 := m.userMap[targetId]
    _, ok := m.userMap[srcId]
    if !ok || !ok2 {
        errChan <- &backendstore.MyStorageError{Message: "Non-exist User ID"}
        return
    }
    srcUser := m.userMap[srcId]
    for e := srcUser.Following.Front(); e != nil; e = e.Next() {
        fuid := e.Value.(uint)
        if fuid == targetId {
            result <- true
            return
        }
    }
    result <- false
}

// User srcId starts to follow targetId.
func (m *MemUserStore) StartFollowingDB(srcId uint, targetId uint, result chan bool,
        errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if targetId == srcId {
        errorChan <- errors.New("cannot follow themselves")
        return
    }
    srcUser, ok2 := m.userMap[srcId]
    targetUser, ok := m.userMap[targetId]
    if !ok || !ok2 {
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist User ID"}
        return
    }
    for e := srcUser.Following.Front(); e != nil; e = e.Next() {
        fuid := e.Value.(uint)
        if fuid == targetId {
            errorChan <- errors.New("already followed")
            return
        }
    }
    newSrcUser := backendstore.UserEntity{}
    newTarUser := backendstore.UserEntity{}
    copyUserEntity(&newSrcUser, srcUser)
    copyUserEntity(&newTarUser, targetUser)
    newSrcUser.Following.PushBack(targetId) // Src Follows Target
    newTarUser.Follower.PushBack(srcId)     // Target Add Follower Src
    m.userMap[srcId] = &newSrcUser
    m.userMap[targetId] = &newTarUser
    result <- true
}

// srcId stop following targetId.
// targetId remove the follower srcId.
func (m *MemUserStore) StopFollowingDB(srcId uint, targetId uint, result chan bool,
        errorChan chan error) {
    m.Lock()
    defer m.Unlock()
    if targetId == srcId {
        errorChan <- errors.New("cannot unfollow yourself")
        return
    }
    srcUser, ok2 := m.userMap[srcId]
    targetUser, ok := m.userMap[targetId]
    if !ok || !ok2 {
        errorChan <- &backendstore.MyStorageError{Message: "Non-exist ID"}
        return
    }
    newSrcUser := &backendstore.UserEntity{}
    newTarUser := &backendstore.UserEntity{}
    copyUserEntity(newSrcUser, srcUser)
    copyUserEntity(newTarUser, targetUser)
    var found bool
    found = false
    for e := newSrcUser.Following.Front(); e != nil; e = e.Next() {
        fuid := e.Value.(uint)
        if fuid == targetId {
            found = true
            newSrcUser.Following.Remove(e) // srcId stop following targetId
        }
    }
    if !found {
        errorChan <- errors.New("src is not following tar")
        return
    }
    for e := newTarUser.Follower.Front(); e != nil; e = e.Next() {
        fuid := e.Value.(uint)
        if fuid == targetId {
            newTarUser.Follower.Remove(e) // targetId remove the follower srcId
        }
    }
    m.userMap[srcId] = newSrcUser
    m.userMap[targetId] = newTarUser
    result <- true
}

// Get the snapshot of userMap and userSet
func (m *MemUserStore) GetSnapshot() ([]byte, error) {
    return m.MarshalJSON()
}

func (m *MemUserStore) MarshalJSON() ([]byte, error) {
    m.Lock()
    defer m.Unlock()
    return json.Marshal(map[string]interface{}{
        "userMap":     m.userMap,
        "userNameSet": m.userNameSet,
        "pkCounter":   m.pkCounter,
    })
}

// Copy a list of uint
func copyUintList(dst *list.List, src *list.List) {
    if dst == nil {
        dst = list.New()
    }
    if src == nil {
        src = list.New()
    }
    for e := src.Front(); e != nil; e = e.Next() {
        pId := e.Value.(uint)
        dst.PushBack(pId)
    }
}

func copyUserEntity(dst *backendstore.UserEntity, src *backendstore.UserEntity) {
    dst.Tweets = list.New()
    dst.Following = list.New()
    dst.Follower = list.New()
    copyUintList(dst.Tweets, src.Tweets)
    copyUintList(dst.Follower, src.Follower)
    copyUintList(dst.Following, src.Following)
    dst.Password = src.Password
    dst.UserName = src.UserName
    dst.ID = src.ID
}
