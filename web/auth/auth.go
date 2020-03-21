package auth

import (
	"log"
	"net/http"
	"zl2501-final-project/web/session"
	_ "zl2501-final-project/web/session/storage/memory"
)

//TODO:
// the functionality of Encrypt and store the private key
// in database(memory)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.GetManagerSingleton("memory")
}

// This is a middleware handler for checking auth
type authMiddleware struct {
	handler http.Handler
}

func (a *authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//session := globalSessions
	ok := globalSessions.SessionAuth(r)
	if ok {
		a.handler.ServeHTTP(w, r)
	} else {
		log.Println("The request is not authenticated. Redirect to index.")
		http.Redirect(w, r, "/", 302) // Go the index
	}
}

// Check weather this request is authenticated
func CheckAuth(handlerToWrap http.Handler) *authMiddleware {
	return &authMiddleware{
		handler: handlerToWrap,
	}
}
