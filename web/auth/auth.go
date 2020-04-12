package auth

import (
	"log"
	"net/http"
)



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
