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
		hash, _ := EncodePassword(password)
		uId, _ := userRepo.CreateNewUser(&repository.UserInfo{
			UserName: userName,
			Password: hash,
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
	if globalSessions.SessionAuth(r) {
		http.Redirect(w, r, "/home", 302)
	} else {
		t, _ := template.ParseFiles("template/index.html")
		t.Execute(w, nil)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	globalSessions.SessionDestroy(w, r)
	http.Redirect(w, r, "/index", 302)
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html")
		w.Header().Set("Content-Type", "text/html")
		_ = t.Execute(w, nil)
	} else {
		userRepo := model.GetUserRepo()
		//userRepo
		_ = r.ParseForm()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		log.Println("username:", r.Form["username"][0])
		log.Println("password:", r.Form["password"][0])
		user := userRepo.SelectByName(userName)
		if user == nil {
			log.Println("User does not exist.")
			http.Redirect(w, r, "/login", 302)
		} else {
			if e := ComparePassword(user.Password, password); e != nil {
				log.Println("Wrong password.")
				http.Redirect(w, r, "/login", 302)
			}
		}
		sess := globalSessions.SessionStart(w, r)
		sess.Set("userName", userName)
		http.Redirect(w, r, "/home", 302)

	}
}
