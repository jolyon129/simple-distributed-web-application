package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
	"zl2501-final-project/web/session/storage"
)

// provides stores the implementation of manager
var provides = make(map[string]ProviderInterface)

// This is a singleton and used across the application.
var GlobalSessionManager *Manager

// Provider interface in order to represent the
// underlying structure of the session manager
type ProviderInterface interface {
	SessionInit(sid string) (storage.SessionStoreInterface, error)
	SessionRead(sid string) (storage.SessionStoreInterface, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

// global session manager
type Manager struct {
	cookieName  string            //private cookiename
	mu          sync.Mutex        // protects session
	provider    ProviderInterface // A bridge to represent the underlying structure of session
	maxlifetime int64
}

func GetManagerSingleton(provideName string) (*Manager, error) {
	if GlobalSessionManager == nil {
		provider, ok := provides[provideName]
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

// Register makes a session manager provider available by the provided name.
func Register(name string, provider ProviderInterface) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

// Generate the unique ID for a session
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// Read sessionId from cookie.
// If not exist, create a new sessionId and inject into cookie.
// If exist, reuse the same session
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session storage.SessionStoreInterface) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" { // If no cookie, a new session
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else { // Read session from cookie
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

// Manually terminate the session and ask clients to overwrite the corresponding cookie
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.mu.Lock()
		defer manager.mu.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

// A background thread to periodically do garbage collection for expired sessions
func (manager *Manager) GC() {
	//manager.mu.Lock()
	//defer manager.mu.Unlock()
	//manager.provider.SessionGC(manager.maxlifetime)
	//time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}
