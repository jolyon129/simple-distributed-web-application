package storage

// Provider interface in order to represent the
// underlying structure of our session manager
type ProviderInterface interface {
	SessionInit(sid string) (SessionStore, error)
	SessionRead(sid string) (SessionStore, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type SessionStore interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}

