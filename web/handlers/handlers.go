package handlers

import (
	"html/template"
	"log"
	"net/http"
	"zl2501-final-project/web/session"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.GetManagerSingleton("memory")
	go globalSessions.GC() // Spawn the garbage collection service when importing
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
