package controller

import (
	"html/template"
	"log"
	"net/http"
	"zl2501-final-project/web/model"
	"zl2501-final-project/web/model/repository"
	"zl2501-final-project/web/session"
	_ "zl2501-final-project/web/session/storage/memory"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.GetManagerSingleton("memory")
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/signup.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		log.Println("username:", r.Form["username"])
		log.Println("password:", r.Form["password"])
		userRepo := model.GetUserRepo()
		uId, _ := userRepo.CreateNewUser(&repository.UserInfo{
			UserName: userName,
			Password: password,
		})
		sess := globalSessions.SessionStart(w, r)
		sess.Set("userName", userName)
		sess.Set("userId", uId)
		http.Redirect(w, r, "/home", 302)
	}
}

type homeView struct {
	Name string
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess := globalSessions.SessionStart(w, r)
		t, _ := template.ParseFiles("template/home.html")
		w.Header().Set("Content-Type", "text/html")
		user := homeView{Name: sess.Get("userName").(string)}
		log.Println(user.Name)
		t.Execute(w, user)
	}
}
func GoIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, nil)
}
