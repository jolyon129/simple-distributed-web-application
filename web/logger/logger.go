package logger

import (
	"log"
	"net/http"
	"os"
	"time"
)

// logRequestsMiddleware is a middleware handler which implement the handler interface
type logRequestsMiddleware struct {
	handler http.Handler
	logger  *log.Logger
}

func (l *logRequestsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	l.logger.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func LogRequests(handlerToWrap http.Handler) *logRequestsMiddleware {
	logger := log.New(os.Stdout, "LogRequests:", log.Ltime|log.Lshortfile)
	return &logRequestsMiddleware{
		handler: handlerToWrap,
		logger:  logger,
	}
}
