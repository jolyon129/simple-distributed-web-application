package controller

import (
	"html/template"
	"log"
	"net/http"
	"zl2501-final-project/web/constant"
	"zl2501-final-project/web/model"
	"zl2501-final-project/web/model/repository"
	"zl2501-final-project/web/session/sessmanager"
)

var globalSessions *sessmanager.Manager

func init() {
	globalSessions, _ = sessmanager.GetManagerSingleton(sessmanager.ProviderName)
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
		uId, error := userRepo.CreateNewUser(&repository.UserInfo{
			UserName: userName,
			Password: hash,
		})
		if error != nil {
			println(error)
			http.Redirect(w, r, "/signup", 302)
		} else {
			sess := globalSessions.SessionStart(w, r)
			sess.Set(constant.UserName, userName)
			sess.Set(constant.UserId, uId)
			http.Redirect(w, r, "/home", 302)
		}
	}
}

type homeView struct {
	Name string
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess := globalSessions.SessionStart(w, r)
		uname := sess.Get(constant.UserName).(string)
		//log.Println("THe user name in session is:",uname)
		t, _ := template.ParseFiles("template/home.html")
		w.Header().Set("Content-Type", "text/html")
		user := homeView{Name: uname}
		log.Println(user.Name)
		t.Execute(w, user)
	}
}

// Go to index page if not logged in.
// If already logged in(session is valid), go to /home instead.
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
	http.Redirect(w, r, "/", 302)
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
		sess.Set(constant.UserName, userName)
		sess.Set(constant.UserId, user.ID)
		http.Redirect(w, r, "/home", 302)

	}
}

func Tweet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/tweet.html")
		w.Header().Set("Content-Type", "text/html")
		_ = t.Execute(w, nil)
	} else {

	}
}
