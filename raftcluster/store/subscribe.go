package store

// This manager implement a observer pattern for listeners
// to subscribe a commandLog. Once the commandLog is committed and
// executed, the listener will be notified.
// Each event should only be notified once.
var managerSingle proposeEventManager

func init() {
    managerSingle = proposeEventManager{proposeListener: make(map[uint64]*value)}
}

// Expose approach to subscribe the event of commit of a propose request
type proposeEventManager struct {
    proposeListener map[uint64]*value
}

type value struct {
    resultC chan interface{}
    errC    chan error
}

// Subscribe to a propose request.
// Listeners can wait on the return result channel and error channel to be waked up.
func (m *proposeEventManager) subscribe(commandID uint64) *value {
    resultC := make(chan interface{})
    errC := make(chan error)
    m.proposeListener[commandID] = &value{
        resultC: resultC,
        errC:    errC,
    }
    return m.proposeListener[commandID]
}

// Notified the listeners if there is someone subscribe to this;
// If not, ignore.
func (m *proposeEventManager) notify(commandID uint64, result interface{}) {
    val, ok := m.proposeListener[commandID]
    if !ok {
        return
    }
    val.resultC <- result
}

// If there is an error, pass error to the error channel and notify the listener
func (m *proposeEventManager) notifyError(cmdID uint64, err error) {
    val, ok := m.proposeListener[cmdID]
    if !ok {
        return
    }
    val.errC <- err
}

func (m *proposeEventManager) unsubscribe(commandID uint64) {
    delete(m.proposeListener, commandID)
}
