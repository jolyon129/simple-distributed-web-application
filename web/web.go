package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"zl2501-final-project/web/auth"
	"zl2501-final-project/web/controller"
	"zl2501-final-project/web/logger"
	"zl2501-final-project/web/session"
	_ "zl2501-final-project/web/session/storage/memory"
)

var globalSessions *session.Manager

// Then, initialize the session manager
func init() {
	// Set global logger
	log.SetPrefix("GlobalLogger: ")
	log.SetFlags(log.Ltime | log.Llongfile)
	log.Println("init started")
	globalSessions, _ = session.GetManagerSingleton("memory")
}

// Adapter Pattern for middleware handlers.
// Ref:
// "https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81"
type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func StartService() {
	//session.Register("memory",nil)
	//globalSessions,_ = session.NewManager("memory","gosessionid",3600)
	//go globalSessions.GC()
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(controller.GoIndex)) // set router
	mux.Handle("/login", http.HandlerFunc(login))
	mux.Handle("/signup", http.HandlerFunc(controller.SignUp))
	mux.Handle("/home", MiddlewareAdapt(http.HandlerFunc(controller.Home), auth.CheckAuth))
	mux.Handle("/logout", http.HandlerFunc(controller.LogOut))

	err := http.ListenAndServe(":9090", logger.LogRequests(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Println("Server starts at: localhost:9090")
	}
}

//type HandlerWrapper func(handler http.Handler) http.Handler

func MiddlewareAdapt(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/home", 302)
	}
}
