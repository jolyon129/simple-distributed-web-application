package controller

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"zl2501-final-project/web/constant"


)

var globalSessions *sessmanager.Manager

func init() {
	globalSessions, _ = sessmanager.GetManagerSingleton(sessmanager.ProviderName)
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"signup.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		if len(userName) == 0 || len(password) == 0 {
			log.Println("Illegal user name or password")
			http.Redirect(w, r, "/signup", 302)
			return
		}
		log.Println("username:", r.Form["username"])
		log.Println("password:", r.Form["password"])
		userRepo := model.GetUserRepo()
		hash, _ := EncodePassword(password)
		uId, error := userRepo.CreateNewUser(&repository.UserInfo{
			UserName: userName,
			Password: hash,
		})
		if error != nil {
			log.Println(error)
			http.Redirect(w, r, "/signup", 302)
		} else {
			sess := globalSessions.SessionStart(w, r)
			sess.Set(constant.UserName, userName)
			sess.Set(constant.UserId, uId)
			http.Redirect(w, r, "/home", 303)
		}
	}
}

// Go to index page if not logged in.
// If already logged in(session is valid), go to /home instead.
func GoIndex(w http.ResponseWriter, r *http.Request) {
	if globalSessions.SessionAuth(r) {
		w.Header().Set("cache-control","no-store") // Avoid the safari remember the redirect
		http.Redirect(w, r, "/home", 302)
	} else {
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"index.html")
		t.Execute(w, nil)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	globalSessions.SessionDestroy(w, r)
	http.Redirect(w, r, "/index", 302)
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"login.html")
		w.Header().Set("Content-Type", "text/html")
		_ = t.Execute(w, nil)
	} else {
		userRepo := model.GetUserRepo()
		//userRepo
		_ = r.ParseForm()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		if len(userName) == 0 || len(password) == 0 {
			log.Println("Illegal user name or password")
			http.Redirect(w, r, "/login", 302)
			return
		}
		log.Println("username:", r.Form["username"][0])
		log.Println("password:", r.Form["password"][0])
		user := userRepo.SelectByName(userName)
		if user == nil {
			log.Println("User does not exist.")
			http.Redirect(w, r, "/login", 302)
			return
		} else {
			if e := ComparePassword(user.Password, password); e != nil {
				log.Println("Wrong password.")
				http.Redirect(w, r, "/login", 302)
				return
			}
		}
		sess := globalSessions.SessionStart(w, r)
		sess.Set(constant.UserName, userName)
		sess.Set(constant.UserId, user.ID)
		http.Redirect(w, r, "/home", 303)

	}
}

func Tweet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"tweet.html")
		w.Header().Set("Content-Type", "text/html")
		_ = t.Execute(w, nil)
	} else {
		r.ParseForm()
		sess := globalSessions.SessionStart(w, r)
		userRepo := model.GetUserRepo()
		postRepo := model.GetPostRepo()
		content := r.Form["content"][0]
		log.Println(content)
		uId := sess.Get(constant.UserId).(uint)
		pId, _ := postRepo.CreateNewPost(repository.PostInfo{
			UserID:  uId,
			Content: content,
		})
		userRepo.AddTweetToUser(uId, pId)
		http.Redirect(w, r, "/home", 302)
	}
}

type user struct {
	Name     string
	Id       string
	Followed bool
}
type viewUserView struct {
	UserList []user
}

// View all users in the system.
// So that the user can follow others
func ViewUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess := globalSessions.SessionStart(w, r)
		myUid := sess.Get(constant.UserId).(uint)
		userRepo := model.GetUserRepo()
		allUsers := userRepo.FindAllUsers()
		newUserList := make([]user, 0)
		for _, value := range allUsers {
			if value.ID == myUid { // Exclude myself
				continue
			}
			newUserList = append(newUserList, user{Name: value.UserName,
				Followed: userRepo.CheckWhetherFollowing(myUid,value.ID),
				Id:       strconv.Itoa(int(value.ID))})
		}
		view := viewUserView{
			UserList: newUserList,
		}
		log.Println(view.UserList)
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"users.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, view)
	}
}

func Follow(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	val, ok := param[constant.UserId]
	if !ok { // If no `userId` in url
		log.Println("No user id to follow!")
	} else {
		uId32, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			log.Println("Wrong user id!")
		}
		uId := uint(uId32)
		sess := globalSessions.SessionStart(w, r)
		myUid := sess.Get(constant.UserId).(uint)
		uRepo := model.GetUserRepo()
		uRepo.StartFollowing(myUid, uId)
		http.Redirect(w, r, "/users", 302)
	}
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	val, ok := param[constant.UserId]
	if !ok { // If no `userId` in url
		log.Println("No user id to unfollow!")
		http.Redirect(w, r, "/users", 302)
	} else {
		uId32, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			log.Println("Wrong user id!")
			http.Redirect(w, r, "/users", 302)
		}
		targetId := uint(uId32)
		sess := globalSessions.SessionStart(w, r)
		myUid := sess.Get(constant.UserId).(uint)
		uRepo := model.GetUserRepo()
		uRepo.StopFollowing(myUid, targetId)
		http.Redirect(w, r, "/users", 302)
	}
}

type userView struct {
	Name     string
	MyTweets []tweet
}

func User(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userNameStr := path.Base(html.UnescapeString(r.URL.Path))
		uE := model.GetUserRepo().SelectByName(userNameStr)
		if uE == nil {
			log.Println("Illegal User Name")
			http.Redirect(w, r, "/users", 302)
			return
		}
		posts := make([]tweet, 0)
		for e := uE.Posts.Back(); e != nil; e = e.Prev() {
			pid := e.Value.(uint)
			p := model.GetPostRepo().SelectById(pid)
			posts = append(posts, tweet{
				Content:   p.Content,
				CreatedAt: p.CreatedTime.Format(constant.TimeFormat),
				CreatedBy: uE.UserName,
				UserId:    int(uE.ID),
			})
		}
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"user.html")
		t.Execute(w, userView{
			Name:     uE.UserName,
			MyTweets: posts,
		})
	}
}
