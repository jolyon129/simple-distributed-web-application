package memory

import (
    "container/list"
    "encoding/json"
    "fmt"
    "sync"
    "time"
    storage "zl2501-final-project/raftcluster/store/authstore"
)

var pder = &Provider{list: list.New()}

// MemSessStore implement the session interface
type MemSessStore struct {
    storage.SessionStorageInterface
    sid          string                      // unique session id
    timeAccessed time.Time                   // last access time
    value        map[interface{}]interface{} // session value stored inside
    sync.Mutex
}

func (st *MemSessStore) Set(key, value interface{}) error {
    st.Lock()
    defer st.Unlock()
    st.value[key] = value
    pder.SessionUpdate(st.sid)
    return nil
}

func (st *MemSessStore) Get(key interface{}) interface{} {
    st.Lock()
    defer st.Unlock()
    pder.SessionUpdate(st.sid)
    if v, ok := st.value[key]; ok {
        return v
    } else {
        return nil
    }
}

func (st *MemSessStore) Delete(key interface{}) error {
    st.Lock()
    defer st.Unlock()
    delete(st.value, key)
    pder.SessionUpdate(st.sid)
    return nil
}

func (st *MemSessStore) SessionID() string {
    st.Lock()
    defer st.Unlock()
    return st.sid
}

// Implement Provider interface.
// Use LRU to store the sessions
type Provider struct {
    storage.ProviderInterface
    lock     sync.Mutex               // lock
    sessions map[string]*list.Element // save in memory
    list     *list.List               // LRU
}

func (pder *Provider) SessionInit(sid string) (storage.SessionStorageInterface, error) {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    v := make(map[interface{}]interface{}, 0)
    newsess := &MemSessStore{sid: sid, timeAccessed: time.Now(), value: v}
    element := pder.list.PushBack(newsess)
    pder.sessions[sid] = element
    return newsess, nil
}

func (pder *Provider) SessionRead(sid string) (storage.SessionStorageInterface, error) {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    if element, ok := pder.sessions[sid]; ok {
        //pder.SessionUpdate(sid)
        return element.Value.(*MemSessStore), nil
    } else {
        return nil, fmt.Errorf("the session Id: %s is not existed", sid)
    }
}

func (pder *Provider) SessionDestroy(sid string) error {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    if element, ok := pder.sessions[sid]; ok {
        delete(pder.sessions, sid)
        pder.list.Remove(element)
        return nil
    }
    return nil
}

// Periodically check the list in Session Store and delete the expired sessions.
func (pder *Provider) SessionGC(maxlifetime int64) {
    pder.lock.Lock()
    defer pder.lock.Unlock()

    for {
        element := pder.list.Back()
        if element == nil {
            break
        }
        if (element.Value.(*MemSessStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
            pder.list.Remove(element)
            delete(pder.sessions, element.Value.(*MemSessStore).sid)
        } else {
            break
        }
    }
}

func (pder *Provider) SessionUpdate(sid string) error {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    if element, ok := pder.sessions[sid]; ok {
        element.Value.(*MemSessStore).timeAccessed = time.Now()
        pder.list.MoveToFront(element)
        return nil
    }
    return nil
}

// Get a snapshot of the data structure under the hood
// Return a json of byte array, like marshall
func (pder *Provider) GetSnapshot() ([]byte, error) {
    pder.lock.Lock()
    defer pder.lock.Unlock()
    // to encode a Go map type it must be of the form map[string]T
    // (where T is any Go type supported by the json package).
    // https://blog.golang.org/json
    sessArr := make([]map[string]interface{}, pder.list.Len())
    i := 0
    for e := pder.list.Front(); e != nil; e = e.Next() {
        sess := e.Value.(*MemSessStore)
        m := make(map[string]interface{})
        for k, v := range sess.value { // convert map[interface]interface to map[string]interface
            m[k.(string)] = v
        }
        m["timeAccessed"] = sess.timeAccessed
        sessArr[i] = m
        i++
    }
    tmp := map[string]interface{}{
        "provider": sessArr,
    }
    return json.Marshal(tmp)
}

func (m *MemSessStore) MarshalJSON() ([]byte, error) {
    value := make(map[string]interface{})
    for k, v := range m.value { // convert map[interface]interface to map[string]interface
        value[k.(string)] = v
    }
    return json.Marshal(map[string]interface{}{
        "timeAccessed": m.timeAccessed,
        "value":        value,
        "sid":          m.sid,
    })
}

func init() {
    pder.sessions = make(map[string]*list.Element, 0)
    storage.RegisterProvider("memory", pder)
}
