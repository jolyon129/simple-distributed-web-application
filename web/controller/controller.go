package controller

import (
    "context"
    "errors"
    "html"
    "html/template"
    "log"
    "net/http"
    "path"
    "strconv"
    "zl2501-final-project/web/constant"
    . "zl2501-final-project/web/pb"
)

type appError struct {
    Err     error
    Message string
    Code    int
}

func (a appError) Error() string {
    return a.Message
}

func SignUp(w http.ResponseWriter, r *http.Request) error {
    if r.Method == "GET" {
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "signup.html")
        w.Header().Set("Content-Type", "text/html")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        userName := r.Form["username"][0]
        password := r.Form["password"][0]
        if len(userName) == 0 || len(password) == 0 {
            http.Redirect(w, r, "/signup", 302)
            return appError{
                Err:     errors.New("illegal user name or password"),
                Message: "Illegal user name or password",
                Code:    400,
            }
        }
        log.Println("username:", r.Form["username"])
        log.Println("password:", r.Form["password"])
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        newUserResponse, error := BackendClientIns.NewUser(ctx, &NewUserRequest{
            UserName: userName,
            UserPwd:  password,
        })
        if error != nil {
            //http.Redirect(w, r, "/signup", 302)
            return appError{
                Err:     error,
                Message: error.Error(),
                Code:    400,
            }
        } else {
            sessId, error := SessionStart(w, r)
            if error != nil {
                return appError{
                    Err:     error,
                    Message: error.Error(),
                    Code:    500,
                }
            }
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
                return appError{
                    Err:     err1,
                    Message: err1.Error(),
                    Code:    500,
                }
            }
            http.Redirect(w, r, "/home", 303)
        }
    }
    return nil
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

func LogOut(w http.ResponseWriter, r *http.Request) error {
    http.Redirect(w, r, "/index", 302)
    return SessionDestroy(w, r)
}

func LogIn(w http.ResponseWriter, r *http.Request) error {
    if r.Method == "GET" {
        t, _ := template.ParseFiles(constant.RelativePathForTemplate + "login.html")
        w.Header().Set("Content-Type", "text/html")
        _ = t.Execute(w, nil)
    } else {
        _ = r.ParseForm()
        userName := r.Form["username"][0]
        password := r.Form["password"][0]
        if len(userName) == 0 || len(password) == 0 {
            return errors.New("Illegal user name or password")
            //http.Redirect(w, r, "/login", 302)
            //return nil
        }
        log.Println("username:", r.Form["username"][0])
        log.Println("password:", r.Form["password"][0])
        //user := userRepo.SelectByName(userName)
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        res, SelectByNameErr := BackendClientIns.UserSelectByName(ctx, &UserSelectByNameRequest{
            Name: userName,
        })
        if SelectByNameErr != nil { // User not existed
            return SelectByNameErr
            //http.Redirect(w, r, "/login", 302)
        }
        if e := ComparePassword(res.User.Password, password); e != nil { // Wrong Password
            //http.Redirect(w, r, "/login", 302)
            return appError{
                Err:     e,
                Message: "Wrong password.",
                Code:    400,
            }
        }
        //ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        sessId, err := SessionInit(w, r)
        if err != nil {
            log.Print(err)
            return appError{
                Err:     err,
                Message: err.Error(),
                Code:    500,
            }
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
            return appError{
                Err:     err1,
                Message: err1.Error(),
                Code:    500,
            }
        }
        http.Redirect(w, r, "/home", 303)
    }
    return nil
}

func Tweet(w http.ResponseWriter, r *http.Request) error {
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
            return err0
        }
        ctx2, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        _, err := BackendClientIns.UserAddTweet(ctx2, &UserAddTweetRequest{
            UserId:  myUid,
            Content: content,
        })
        if err != nil {
            log.Print(err)
            return err
        }
        http.Redirect(w, r, "/home", 302)
    }
    return nil
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
func ViewUsers(w http.ResponseWriter, r *http.Request) error {
    if r.Method == "GET" {
        myUid, err0 := GetMyUserId(r)
        if err0 != nil {
            return err0
        }
        ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        response, err := BackendClientIns.FindAllUsers(ctx1, &FindAllUsersRequest{})
        if err != nil {
            return err
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
    return  nil
}

func Follow(w http.ResponseWriter, r *http.Request) error{
    param := r.URL.Query()
    val, ok := param[constant.UserId]
    if !ok { // If no `userId` in url
        log.Println("No user id to follow!")
    } else {
        targetUid, err := strconv.ParseUint(val[0], 10, 64)
        if err != nil {
            return errors.New("Wrong user id")
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
    return nil
}

func Unfollow(w http.ResponseWriter, r *http.Request) error {
    param := r.URL.Query()
    val, ok := param[constant.UserId]
    if !ok { // If no `userId` in url
        log.Println("No user id to unfollow!")
        http.Redirect(w, r, "/users", 302)
    } else {
        targetUid, err := strconv.ParseUint(val[0], 10, 64)
        if err != nil {
            errors.New("Wrong user id!")
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
    return nil
}

type userView struct {
    Name     string
    MyTweets []tweet
}

func User(w http.ResponseWriter, r *http.Request) error {
    if r.Method == "GET" {
        userNameStr := path.Base(html.UnescapeString(r.URL.Path))
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        res, _ := BackendClientIns.UserSelectByName(ctx, &UserSelectByNameRequest{
            Name: userNameStr,
        })
        userE := res.User
        if userE == nil {
            return errors.New("Illegal User Name")
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
    return nil
}
