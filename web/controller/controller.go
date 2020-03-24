package controller

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
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
		http.Redirect(w, r, "/home", 302)

	}
}

func Tweet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/tweet.html")
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
		myUserE := userRepo.SelectById(myUid)
		myFollowingMap := make(map[uint]bool)
		for e := myUserE.Following.Front(); e != nil; e = e.Next() {
			myFollowingMap[e.Value.(uint)] = true
		}
		for _, value := range allUsers {
			if value.ID == myUid { // Exclude myself
				continue
			}
			if _, ok := myFollowingMap[value.ID]; ok {
				newUserList = append(newUserList, user{Name: value.UserName,
					Followed: true,
					Id:       strconv.Itoa(int(value.ID))})
			} else {
				newUserList = append(newUserList, user{Name: value.UserName,
					Followed: false,
					Id:       strconv.Itoa(int(value.ID))})
			}
		}
		view := viewUserView{
			UserList: newUserList,
		}
		log.Println(view.UserList)
		t, _ := template.ParseFiles("template/users.html")
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
		uidstr := path.Base(html.UnescapeString(r.URL.Path))
		uid, err := strconv.ParseUint(uidstr, 10, 32)
		if err != nil {
			log.Println("Illegal userid")
			http.Redirect(w, r, "/users", 302)
			return
		}
		uiduint := uint(uid)
		uE := model.GetUserRepo().SelectById(uiduint)
		if uE == nil {
			log.Println("Illegal userid")
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
		t, _ := template.ParseFiles("template/user.html")
		t.Execute(w, userView{
			Name:     uE.UserName,
			MyTweets: posts,
		})
	}
}
