package web

import (
	"log"
	"net/http"
	"zl2501-final-project/web/auth"
	"zl2501-final-project/web/constant"
	"zl2501-final-project/web/controller"
	"zl2501-final-project/web/logger"
	"zl2501-final-project/web/session/sessmanager"
)

var globalSessions *sessmanager.Manager

// Then, initialize the session manager
func init() {
	// Set global logger
	log.SetPrefix("GlobalLogger: ")
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.Println("init started")
	globalSessions, _ = sessmanager.GetManagerSingleton("memory")
}

func StartService() {
	//session.RegisterProvider("memory",nil)
	//globalSessions,_ = session.NewManager("memory","gosessionid",3600)
	//go globalSessions.GC()
	mux := http.NewServeMux()
	mux.Handle("/", MiddlewareAdapt(http.HandlerFunc(controller.GoIndex), SetHeader))                                  // set router
	mux.Handle("/index", MiddlewareAdapt(http.HandlerFunc(controller.GoIndex), SetHeader)) // set router
	mux.Handle("/login", MiddlewareAdapt(http.HandlerFunc(controller.LogIn), SetHeader))
	mux.Handle("/signup", MiddlewareAdapt(http.HandlerFunc(controller.SignUp), SetHeader))
	mux.Handle("/home", MiddlewareAdapt(http.HandlerFunc(controller.Home), auth.CheckAuth, SetHeader))
	mux.Handle("/logout", MiddlewareAdapt(http.HandlerFunc(controller.LogOut), SetHeader))
	mux.Handle("/tweet", MiddlewareAdapt(http.HandlerFunc(controller.Tweet), auth.CheckAuth, SetHeader))
	mux.Handle("/users", MiddlewareAdapt(http.HandlerFunc(controller.ViewUsers), auth.CheckAuth, SetHeader))
	mux.Handle("/user/", MiddlewareAdapt(http.HandlerFunc(controller.User), auth.CheckAuth, SetHeader))
	mux.Handle("/follow", MiddlewareAdapt(http.HandlerFunc(controller.Follow), auth.CheckAuth, SetHeader))
	mux.Handle("/unfollow", MiddlewareAdapt(http.HandlerFunc(controller.Unfollow), auth.CheckAuth, SetHeader))
	err := http.ListenAndServe(":"+constant.Port, logger.LogRequests(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Println("Server starts at: localhost:"+constant.Port)
	}
}

// Adapt all middlewares to the handler.
// The function will call them one by one (in reverse order) in a chained manner,
// returning the result of the first adapter.
// Ref: https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
func MiddlewareAdapt(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

// This is a middleware
// Add Some Header
func SetHeader(handlerToWrap http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "text/html")
		w.Header().Set("cache-control", "no-store")
		handlerToWrap.ServeHTTP(w, r)
	})
}
