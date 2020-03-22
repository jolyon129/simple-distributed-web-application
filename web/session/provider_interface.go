package session

import (
	"zl2501-final-project/web/session/storage"
)

// provides stores the implementation of manager
var provides = make(map[string]ProviderInterface)

// Provider interface in order to represent the
// underlying structure of the session manager
type ProviderInterface interface {
	SessionInit(sid string) (storage.SessionStorageInterface, error)
	// Read session through ssid.
	// If not existed, return (nil, error)
	SessionRead(sid string) (storage.SessionStorageInterface, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

// global session manager
//type Manager struct {
//	cookieName  string            //private cookiename
//	mu          sync.Mutex        // protects session
//	provider    ProviderInterface // A bridge to represent the underlying structure of session
//	maxlifetime int64
//}
//
//func GetManagerSingleton(provideName string) (*Manager, error) {
//	if GlobalSessionManager == nil {
//		provider, ok := provides[provideName]
//		if !ok {
//			return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
//		}
//		GlobalSessionManager = &Manager{provider: provider, cookieName: CookieName, maxlifetime: MaxLifeTime}
//		// Spawn another thread for garbage collection
//		go GlobalSessionManager.GC()
//		return GlobalSessionManager, nil
//	} else {
//		return GlobalSessionManager, nil
//	}
//}

// RegisterProvider makes a session manager provider available by the provided name.
func RegisterProvider(name string, provider ProviderInterface) {
	if provider == nil {
		panic("session: RegisterProvider provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: RegisterProvider called twice for provider " + name)
	}
	provides[name] = provider
}

// Get provider of the implementation of session by name
func GetProvider(name string) (ProviderInterface, bool) {
	provider, ok := provides[name]
	return provider, ok
}
