package store

import "sync"

var managerSingle proposeEventManager

func init() {
    managerSingle = proposeEventManager{proposeListener: make(map[string]*sync.WaitGroup)}
}

// Expose approach to subscribe the event of commit of a propose request
type proposeEventManager struct {
    proposeListener map[string]*sync.WaitGroup
}

// Singleton
func GetProposeEventManager() *proposeEventManager {
    return &managerSingle
}

// Subscribe to the propose request.
// Users can use the return waitgroup to wait for the commit.
func (m *proposeEventManager) subscribe(proposeId string) *sync.WaitGroup {
    var wg sync.WaitGroup
    wg.Add(1)
    m.proposeListener[proposeId] = &wg
    return &wg
}

// notified the listeners
func (m *proposeEventManager) notify(proposeId string) {
    wg := m.proposeListener[proposeId]
    wg.Done()
}
