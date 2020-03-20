package storage

// Provider interface in order to represent the
// underlying structure of our session manager
type ProviderInterface interface {
	SessionInit(sid string) (SessionStoreInterface, error)
	SessionRead(sid string) (SessionStoreInterface, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type SessionStoreInterface interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}
