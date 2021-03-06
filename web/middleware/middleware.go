package middleware

import (
    "context"
    "errors"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"
    "zl2501-final-project/web/constant"
    "zl2501-final-project/web/pb"
)

func init() {
}

type logRequestsMiddleware struct {
    handler http.Handler
    logger  *log.Logger
}

func (l *logRequestsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    l.handler.ServeHTTP(w, r)
    l.logger.Printf("Request:%s %s, Time: %v", r.Method, r.URL.Path, time.Since(start))
}

// logRequestsMiddleware is a middleware handler which implement the handler interface
func LogRequests(handlerToWrap http.Handler) *logRequestsMiddleware {
    logger := log.New(os.Stdout, "LogRequests:", log.Ltime|log.Lshortfile)
    return &logRequestsMiddleware{
        handler: handlerToWrap,
        logger:  logger,
    }
}

// This middleware helps consume the returned error from custom handler!
type AppHandler func(http.ResponseWriter, *http.Request) error

type errorView struct {
    Message string
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := fn(w, r); err != nil {
        //http.Error(w, err.Error(), 500)
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "error.html")
        t.Execute(w, errorView{Message: err.Error()})
    }
}

// This is a middleware handler used to check weather this request is authenticated.
// If not, redirect to the index.
func CheckAuth(handlerToWrap http.Handler) http.Handler {
    return AppHandler(func(w http.ResponseWriter, r *http.Request) error {
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        cookie, err := r.Cookie(constant.SessCookieName)
        if err != nil || cookie.Value == "" {
            log.Printf("Request:%s %s is not authenticated. Redirect to index.", r.Method, r.URL.Path)
            http.Redirect(w, r, "/", 307) // Go the index
        } else {
            sid, _ := url.QueryUnescape(cookie.Value)
            response, err := pb.AuthClientIns.SessionAuth(ctx, &pb.SessionGeneralRequest{
                SessionId: sid,
            })
            if err != nil {
                log.Printf(err.Error())
                //http.Redirect(w, r, "/", 307) // Go the index
                return err
            }
            if response.Ok {
                log.Printf("Request:%s %s is authenticated. SessId: %s", r.Method,
                    r.URL.Path, sid)
                handlerToWrap.ServeHTTP(w, r)
            }
        }
        return nil
    })
}

// This is a middleware to
// add Some Header to response
func SetHeader(handlerToWrap http.Handler) http.Handler {
    return AppHandler(func(w http.ResponseWriter, r *http.Request) error {
        if w == nil {
            handlerToWrap.ServeHTTP(w, r)
            return errors.New("something went wrong")
        }
        w.Header().Set("Content-Type", "text/html")
        w.Header().Set("cache-control", "no-store")
        handlerToWrap.ServeHTTP(w, r)
        return nil
    })
}
