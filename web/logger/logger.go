// Package logger provides the LogRequests middleware.
package logger

import (
	"log"
	"net/http"
	"os"
	"time"
	"zl2501-final-project/auth/sessmanager"
	"zl2501-final-project/web/constant"
)

var globalSessions *sessmanager.Manager

func init() {
	globalSessions, _ = sessmanager.GetManagerSingleton("memory")
}

// logRequestsMiddleware is a middleware handler which implement the handler interface
type logRequestsMiddleware struct {
	handler http.Handler
	logger  *log.Logger
}

func (l *logRequestsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	if globalSessions.SessionAuth(r) {
		session := globalSessions.SessionStart(w, r)
		uname := session.Get(constant.UserName)
		l.logger.Printf("Request:%s %s, Time: %v, User:%s", r.Method, r.URL.Path, time.Since(start), uname)
	} else {
		l.logger.Printf("Request:%s %s, Time: %v, User not logged in", r.Method, r.URL.Path, time.Since(start))
	}

}

func LogRequests(handlerToWrap http.Handler) *logRequestsMiddleware {
	logger := log.New(os.Stdout, "LogRequests:", log.Ltime|log.Lshortfile)
	return &logRequestsMiddleware{
		handler: handlerToWrap,
		logger:  logger,
	}
}
