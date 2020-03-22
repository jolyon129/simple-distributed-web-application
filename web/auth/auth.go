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

// This is a middleware handler used to check weather this request is authenticated.
// If not, redirect to the index.
func CheckAuth(handlerToWrap http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := globalSessions.SessionAuth(r)
		if ok {
			handlerToWrap.ServeHTTP(w, r)
		} else {
			log.Printf("Request:%s %s is not authenticated. Redirect to index.", r.Method, r.URL.Path)
			http.Redirect(w, r, "/", 302) // Go the index
		}
	})
}
