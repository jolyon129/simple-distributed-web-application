// Package manager provides the session manager to manage all sessions.
package sessmanager

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "sync"
    "time"
    "zl2501-final-project/auth/storage"
    _ "zl2501-final-project/auth/storage/raftclient" // Use raft wrapper implementation of session
)

// This is a singleton and used across the application.
var GlobalSessionManager *Manager

// global session manager
type Manager struct {
    cookieName  string                    //private cookiename
    mu          sync.Mutex                // protects session
    provider    storage.ProviderInterface // A bridge to represent the underlying structure of session
    maxlifetime int64
}

// Get the singleton of manager
func GetManagerSingleton(provideName string) (*Manager, error) {
    if GlobalSessionManager == nil {
        provider, ok := storage.GetProvider("raft")
        if !ok {
            return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
        }
        GlobalSessionManager = &Manager{provider: provider, cookieName: CookieName, maxlifetime: MaxLifeTime}
        // Spawn another thread for garbage collection
        go GlobalSessionManager.GC()
        GlobalSessionManager.provider.SessionGC(GlobalSessionManager.maxlifetime)
        return GlobalSessionManager, nil
    } else {
        return GlobalSessionManager, nil
    }
}

// Generate the unique ID for a session
func (manager *Manager) newSessionId() string {
    b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

// Read session from given sessId If its legal.
// If not exist, create a new newSessionId and return.
// If exist and the newSessionId is valid, reuse the same session and return the same one.
func (manager *Manager) SessionStart(ctx context.Context, sessId string) (string, error) {
    result := make(chan storage.SessionStorageInterface)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    manager.mu.Lock()
    defer manager.mu.Unlock()
    go func() {
        if sessId == "" { // Empty
            newSessId := manager.newSessionId()
            sess, err := manager.provider.SessionInit(newSessId)
            if err != nil {
                errorChan <- err
                return
            } else {
                result <- sess
                return
            }
        } else { // sessId is not legal
            oldSess, err := manager.provider.SessionRead(sessId)
            if err != nil {
                newSessId := manager.newSessionId()
                sess, err := manager.provider.SessionInit(newSessId)
                if err != nil {
                    errorChan <- err
                    return
                } else {
                    result <- sess
                    return
                }
            } else { // If the sessId is legal, reuse the session
                if err != nil {
                    errorChan <- err
                    return
                } else {
                    result <- oldSess
                    return
                }
            }
        }
    }()
    select {
    case err := <-errorChan:
        return "", err
    case sess := <-result:
        return sess.SessionID(), nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

// Check whether the session has expired
func (manager *Manager) SessionAuth(ctx context.Context, sessId string) (bool, error) {
    errorChan := make(chan error)
    resultChan := make(chan bool)
    defer close(resultChan)
    defer close(errorChan)
    go func() {
        if _, err := manager.provider.SessionRead(sessId); err != nil {
            errorChan <- err
            return
        } else {
            resultChan <- true
            return
        }
    }()
    select {
    case <-ctx.Done():
        return false, ctx.Err()
    case err := <-errorChan:
        return false, err
    case ret := <-resultChan:
        return ret, nil
    }

}

// Manually terminate the session and ask clients to overwrite the corresponding cookie
func (manager *Manager) SessionDestroy(ctx context.Context, sessId string) (bool, error) {
    errorChan := make(chan error)
    result := make(chan bool)
    defer close(result)
    defer close(errorChan)
    go func() {
        err := manager.provider.SessionDestroy(sessId)
        if err != nil {
            errorChan <- err
            return
        } else {
            result <- true
            return
        }
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// A background thread to periodically do garbage collection for expired sessions
func (manager *Manager) GC() {
    manager.mu.Lock()
    manager.provider.SessionGC(manager.maxlifetime)
    manager.mu.Unlock()
    time.AfterFunc(30*time.Second, func() { manager.GC() })
}

// Set Key and Value to a Session
func (manager *Manager) SetValue(ctx context.Context, sessId string, key,
        value interface{}) (bool, error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        session, err := manager.provider.SessionRead(sessId)
        if err != nil {
            errorChan <- err
            return
        }
        err1 := session.Set(key, value)
        if err1 != nil {
            errorChan <- err
            return
        }
        result <- true
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

// Delete Key and Value from a Session
func (manager *Manager) DeleteValue(ctx context.Context, sessId string, key interface{}) (bool,
        error) {
    result := make(chan bool)
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        sess, err := manager.provider.SessionRead(sessId)
        if err != nil {
            errorChan <- err
            return
        }
        err1 := sess.Delete(key)
        if err1 != nil {
            errorChan <- err
            return
        }
        result <- true
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}

//  Get Value from a Session with the Key
func (manager *Manager) GetValue(ctx context.Context, sessId string, key interface{}) (interface{}, error) {
    result := make(chan interface{})
    errorChan := make(chan error)
    defer close(result)
    defer close(errorChan)
    go func() {
        session, err := manager.provider.SessionRead(sessId)
        if err != nil {
            errorChan <- err
            return
        }
        value := session.Get(key)
        result <- value
    }()
    select {
    case ret := <-result:
        return ret, nil
    case err := <-errorChan:
        return false, err
    case <-ctx.Done():
        return false, ctx.Err()
    }
}
