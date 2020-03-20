package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"zl2501-final-project/web/session"
	_ "zl2501-final-project/web/session/storage/memory"
)

var globalSessions *session.Manager

// Then, initialize the session manager
func init() {
	// Set logger
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ltime | log.Llongfile)
	log.Println("init started")
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC() // Spawn the garbage collection service when importing
}

func StartService() {
	//session.Register("memory",nil)
	//globalSessions,_ = session.NewManager("memory","gosessionid",3600)
	//go globalSessions.GC()
	http.HandleFunc("/", sayHello) // set router
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/home", home)
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Println("Server starts at: localhost:9090")
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, nil)
}

type homeUser struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess := globalSessions.SessionStart(w, r)
		t, _ := template.ParseFiles("template/home.html")
		w.Header().Set("Content-Type", "text/html")
		user := homeUser{Name: sess.Get("username").([]string)[0]}
		log.Println(user.Name)
		t.Execute(w, user)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/signup.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		log.Println("username:", r.Form["username"])
		log.Println("password:", r.Form["password"])
		sess := globalSessions.SessionStart(w, r)
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/home", 302)
	}
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
