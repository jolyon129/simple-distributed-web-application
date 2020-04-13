package controller

import (
    "context"
    "html"
    "html/template"
    "log"
    "net/http"
    "path"
    "strconv"
    "zl2501-final-project/web/constant"
    . "zl2501-final-project/web/pb"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "signup.html")
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
        hash, _ := EncodePassword(password)
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        newUserResponse, error := BackendClientIns.NewUser(ctx, &NewUserRequest{
            UserName: userName,
            UserPwd:  hash,
        })
        if error != nil {
            log.Println(error)
            http.Redirect(w, r, "/signup", 302)
        } else {
            sessId, _ := SessionStart(w, r)
            _, err1 := AuthClientIns.SetValue(ctx, &SetValueRequest{
                Key:   constant.UserName,
                Value: userName,
                Ssid:  sessId,
            })
            _, err2 := AuthClientIns.SetValue(ctx, &SetValueRequest{
                Key:   constant.UserId,
                Value: strconv.Itoa(int(newUserResponse.UserId)),
                Ssid:  sessId,
            })
            if err1 != nil || err2 != nil {
                log.Print(err1)
                log.Print(err2)
                //Todo: Error Handler
            }
            http.Redirect(w, r, "/home", 303)
        }
    }
}

// Go to index page if not logged in.
// If already logged in(session is valid), go to /home instead.
func GoIndex(w http.ResponseWriter, r *http.Request) {
    if CheckAuthRequest(r) {
        w.Header().Set("cache-control", "no-store") // Avoid the safari remember the redirect
        http.Redirect(w, r, "/home", 302)
    } else {
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "index.html")
        t.Execute(w, nil)
    }
}

func LogOut(w http.ResponseWriter, r *http.Request) {
    SessionDestroy(w, r)
    http.Redirect(w, r, "/index", 302)
}

func LogIn(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "login.html")
        w.Header().Set("Content-Type", "text/html")
        _ = t.Execute(w, nil)
    } else {
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
        //user := userRepo.SelectByName(userName)
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        res, SelectByNameErr := BackendClientIns.UserSelectByName(ctx, &UserSelectByNameRequest{
            Name: userName,
        })
        if SelectByNameErr != nil {
            log.Printf(SelectByNameErr.Error())
        }
        //log.Print(err)
        if res.User == nil {
            log.Println("User does not exist.")
            http.Redirect(w, r, "/login", 302)
            return
        } else {
            if e := ComparePassword(res.User.Password, password); e != nil {
                log.Println("Wrong password.")
                http.Redirect(w, r, "/login", 302)
                return
            }
        }
        //ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        sessId, err := SessionInit(w, r)
        if err != nil {
            log.Print(err)
            return //TODO: Add error handler
        }
        _, err1 := AuthClientIns.SetValue(ctx, &SetValueRequest{
            Key:   constant.UserName,
            Value: userName,
            Ssid:  sessId,
        })
        _, err2 := AuthClientIns.SetValue(ctx, &SetValueRequest{
            Key:   constant.UserId,
            Value: strconv.Itoa(int(res.User.UserId)),
            Ssid:  sessId,
        })
        if err1 != nil || err2 != nil {
            log.Print(err1)
            log.Print(err2)
            //Todo: Error Handler
        }
        http.Redirect(w, r, "/home", 303)

    }
}

func Tweet(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "tweet.html")
        w.Header().Set("Content-Type", "text/html")
        _ = t.Execute(w, nil)
    } else {
        r.ParseForm()
        content := r.Form["content"][0]
        log.Println(content)
        myUid, err0 := GetMyUserId(r)
        if err0 != nil {
            log.Print(err0)
        }
        ctx2, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        _, err := BackendClientIns.UserAddTweet(ctx2, &UserAddTweetRequest{
            UserId:  myUid,
            Content: content,
        })
        if err != nil {
            log.Print(err)
        }
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
        myUid, err0 := GetMyUserId(r)
        if err0 != nil {
            log.Print(err0)
        }
        ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        response, err := BackendClientIns.FindAllUsers(ctx1, &FindAllUsersRequest{})
        if err != nil {
            log.Print(err)
        }
        allUsers := response.Users
        newUserList := make([]user, 0)
        for _, value := range allUsers {
            if value.UserId == myUid { // Exclude myself
                continue
            }
            ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
            responseFromWhetherFollowing, _ := BackendClientIns.UserCheckWhetherFollowing(ctx,
                &UserCheckWhetherFollowingRequest{
                    SourceUserId: myUid,
                    TargetUserId: value.UserId,
                })
            newUserList = append(newUserList, user{Name: value.UserName,
                Followed: responseFromWhetherFollowing.Ok,
                Id:       strconv.Itoa(int(value.UserId))})
        }
        view := viewUserView{
            UserList: newUserList,
        }
        log.Println(view.UserList)
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "users.html")
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
        targetUid, err := strconv.ParseUint(val[0], 10, 64)
        if err != nil {
            log.Println("Wrong user id!")
        }
        myUid, _ := GetMyUserId(r)
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        //uRepo.StartFollowing(myUid, uId)
        BackendClientIns.StartFollowing(ctx, &StartFollowingRequest{
            SourceUserId: myUid,
            TargetUserId: targetUid,
        })
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
        targetUid, err := strconv.ParseUint(val[0], 10, 64)
        if err != nil {
            log.Println("Wrong user id!")
        }
        myUid, _ := GetMyUserId(r)
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        //uRepo.StartFollowing(myUid, uId)
        BackendClientIns.StopFollowing(ctx, &StopFollowingRequest{
            SourceUserId: myUid,
            TargetUserId: targetUid,
        })
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
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        res, _ := BackendClientIns.UserSelectByName(ctx, &UserSelectByNameRequest{
            Name: userNameStr,
        })
        //uE := model.GetUserRepo().SelectByName(userNameStr)
        userE := res.User
        if userE == nil {
            log.Println("Illegal User Name")
            http.Redirect(w, r, "/users", 302)
            return
        }
        // Iterate in reverse order because the latest one is stored in the tail in DB
        tweets := make([]tweet, 0)
        for i := len(userE.Tweets) - 1; i >= 0; i-- {
            ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
            res, _ := BackendClientIns.TweetSelectById(ctx, &SelectByIdRequest{
                Id: userE.Tweets[i],
            })
            tweets = append(tweets, tweet{
                Content:   res.Msg.Content,
                CreatedAt: res.Msg.CreatedTime,
                CreatedBy: userE.UserName,
                UserId:    int(userE.UserId),
            })
        }
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "user.html")
        t.Execute(w, userView{
            Name:     userE.UserName,
            MyTweets: tweets,
        })
    }
}
