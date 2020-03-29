package web

import (
	"log"
	"net/http"
	"zl2501-final-project/web/auth"
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
	mux.Handle("/", http.HandlerFunc(controller.GoIndex)) // set router
	mux.Handle("/login", http.HandlerFunc(controller.LogIn))
	mux.Handle("/signup", http.HandlerFunc(controller.SignUp))
	mux.Handle("/home", MiddlewareAdapt(http.HandlerFunc(controller.Home), auth.CheckAuth))
	mux.Handle("/logout", http.HandlerFunc(controller.LogOut))
	mux.Handle("/tweet", MiddlewareAdapt(http.HandlerFunc(controller.Tweet), auth.CheckAuth))
	mux.Handle("/users", MiddlewareAdapt(http.HandlerFunc(controller.ViewUsers), auth.CheckAuth))
	mux.Handle("/user/", MiddlewareAdapt(http.HandlerFunc(controller.User), auth.CheckAuth))
	mux.Handle("/follow", MiddlewareAdapt(http.HandlerFunc(controller.Follow), auth.CheckAuth))
	mux.Handle("/unfollow", MiddlewareAdapt(http.HandlerFunc(controller.Unfollow), auth.CheckAuth))
	err := http.ListenAndServe(":9090", logger.LogRequests(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Println("Server starts at: localhost:9090")
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
