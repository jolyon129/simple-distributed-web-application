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
    "zl2501-final-project/backend/session"
    "zl2501-final-project/backend/session/storage"
    _ "zl2501-final-project/backend/session/storage/memory" // Use memory implementation of session
)

// This is a singleton and used across the application.
var GlobalSessionManager *Manager

// global session manager
type Manager struct {
    cookieName  string                    //private cookiename
    mu          sync.Mutex                // protects session
    provider    session.ProviderInterface // A bridge to represent the underlying structure of session
    maxlifetime int64
}

// Get the singleton of manager
func GetManagerSingleton(provideName string) (*Manager, error) {
    if GlobalSessionManager == nil {
        provider, ok := session.GetProvider("memory")
        if !ok {
            return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
        }
        GlobalSessionManager = &Manager{provider: provider, cookieName: CookieName, maxlifetime: MaxLifeTime}
        // Spawn another thread for garbage collection
        go GlobalSessionManager.GC()
        return GlobalSessionManager, nil
    } else {
        return GlobalSessionManager, nil
    }
}

// Generate the unique ID for a session
func (manager *Manager) sessionId() string {
    b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

// Read sessionId from Given sessId If its legal.
// If not exist, create a new sessionId and return.
// If exist and the sessionId is valid, reuse the same session and return the same one.
func (manager *Manager) SessionStart(ctx context.Context, sessId string) (string, error) {
    result := make(chan storage.SessionStorageInterface)
    errorChan := make(chan error)
    manager.mu.Lock()
    defer manager.mu.Unlock()
    go func() {
        if sessId == "" {
            newSessId := manager.sessionId()
            sess, err := manager.provider.SessionInit(newSessId)
            if err != nil {
                errorChan <- err
            } else {
                result <- sess
            }
        } else {
            oldSess, err := manager.provider.SessionRead(sessId)
            if err != nil {
                newSessId := manager.sessionId()
                sess, err := manager.provider.SessionInit(newSessId)
                if err != nil {
                    errorChan <- err
                } else {
                    result <- sess
                }
            } else { // If the sessId is legal, reuse the session
                if err != nil {
                    errorChan <- err
                } else {
                    result <- oldSess
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
    go func() {
        if _, err := manager.provider.SessionRead(sessId); err != nil {
            errorChan <- err
        } else {
            resultChan <- true
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
func (manager *Manager) SessionDestroy(ctx context.Context, sessId string) (bool,error) {
    errorChan := make(chan error)
    result := make(chan bool)
    go func() {
        err := manager.provider.SessionDestroy(sessId)
        if err != nil {
            errorChan <- err
        } else {
            result <- true
        }
    }()
    select {
    case ret:=<-result:
        return ret,nil
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
    time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}
