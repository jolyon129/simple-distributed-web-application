package authstorage

// provides stores the implementation of manager
var provides = make(map[string]ProviderInterface)

// Provider interface in order to represent the
// underlying structure of the session manager
type ProviderInterface interface {
	// Create a new session
	SessionInit(sid string) (SessionStorageInterface, error)
	// Read session through ssid.
	// If not existed, return (nil, error)
	SessionRead(sid string) (SessionStorageInterface, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
	GetSnapshot()([]byte,error)
}

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
