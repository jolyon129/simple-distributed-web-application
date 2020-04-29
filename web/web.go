package web

import (
    "log"
    "net/http"
    "zl2501-final-project/web/constant"
    "zl2501-final-project/web/controller"
    . "zl2501-final-project/web/middleware"
    "zl2501-final-project/web/pb"
)

// Then, initialize the session manager
func init() {
    // Set global logger
    log.SetPrefix("GlobalLogger: ")
    log.SetFlags(log.Ltime | log.Lshortfile)
    // Set up a connection to the server.

}

// This middleware helps consume the returned error from custom handler!
type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := fn(w, r); err != nil {
        http.Error(w, err.Error(), 500)
    }
}

func StartService() {
    authConn := pb.CreateAuthServiceConnection()
    backendConn := pb.CreateBackendServiceConnection()
    defer authConn.Close()
    defer backendConn.Close()
    mux := http.NewServeMux()
    mux.Handle("/", MiddlewareAdapt(http.HandlerFunc(controller.GoIndex),
        SetHeader)) // set router
    mux.Handle("/index", MiddlewareAdapt(http.HandlerFunc(controller.GoIndex), SetHeader))
    mux.Handle("/login", MiddlewareAdapt(appHandler(controller.LogIn), SetHeader))
    mux.Handle("/signup", MiddlewareAdapt(appHandler(controller.SignUp), SetHeader))
    mux.Handle("/home", MiddlewareAdapt(appHandler(controller.Home), CheckAuth,
        SetHeader))
    mux.Handle("/logout", MiddlewareAdapt(http.HandlerFunc(controller.LogOut), SetHeader))
    mux.Handle("/tweet", MiddlewareAdapt(http.HandlerFunc(controller.Tweet),
        CheckAuth, SetHeader))
    mux.Handle("/users", MiddlewareAdapt(http.HandlerFunc(controller.ViewUsers),
        CheckAuth, SetHeader))
    mux.Handle("/user/", MiddlewareAdapt(http.HandlerFunc(controller.User), CheckAuth,
        SetHeader))
    mux.Handle("/follow", MiddlewareAdapt(http.HandlerFunc(controller.Follow),
        CheckAuth, SetHeader))
    mux.Handle("/unfollow", MiddlewareAdapt(http.HandlerFunc(controller.Unfollow),
        CheckAuth, SetHeader))
    log.Println("Server is going to start at: http://localhost:" + constant.Port)
    log.Fatal(http.ListenAndServe(":"+constant.Port, LogRequests(mux)))
}

// Adapt all middleware to the handler.
// The function will call them one by one (in reverse order) in a chained manner,
// returning the result of the first adapter.
// Ref: https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
func MiddlewareAdapt(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
    for _, mw := range middleware {
        h = mw(h)
    }
    return h
}

