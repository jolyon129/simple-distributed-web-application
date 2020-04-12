package memory

import (
	"container/list"
	"fmt"
	"sync"
	"time"
	"zl2501-final-project/backend/session"
	"zl2501-final-project/backend/session/storage"
)

var pder = &Provider{list: list.New()}

// MemSessStore implement the session interface
type MemSessStore struct {
	session.ProviderInterface
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

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.RegisterProvider("memory", pder)
}
