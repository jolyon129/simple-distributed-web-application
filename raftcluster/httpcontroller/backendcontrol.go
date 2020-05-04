package httpcontroller

import (
    "encoding/json"
    "net/http"
    "strconv"
    . "zl2501-final-project/raftcluster/store"
    "zl2501-final-project/raftcluster/store/backendstore"
)

// Get the user info from request.from
func UserCreate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    //for k, v := range r.Form {
    //    println(k)
    //    println(v)
    //}
    userName := r.Form["username"][0]
    password := r.Form["password"][0]
    uid, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserCreate,
        UserCreateParams{UserInfo{
            UserName: userName,
            Password: password,
        }})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    nres := map[string]uint{
        "uid": uid.(uint),
    }
    ret, _ := json.Marshal(nres)
    w.Write(ret)
}

func convertStrToUint(s string) uint {
    myUidint, _ := strconv.Atoi(s)
    return uint(myUidint)
}

//func UserDelete(w http.ResponseWriter, r *http.Request)                  {}
func UserRead(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    uid := getRouteParam(r, "uid")
    nUid := convertStrToUint(uid)
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserGet,
        UserIDParams{nUid})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusCreated)
    user := res.(*backendstore.UserEntity)
    ret, _ := json.Marshal(user)
    w.Write(ret)
}
func UserUpdate(w http.ResponseWriter, r *http.Request)                  {}
func UserFindAll(w http.ResponseWriter, r *http.Request)                 {}
func UserAddTweetToUserDB(w http.ResponseWriter, r *http.Request)        {}
func UserCheckWhetherFollowingDB(w http.ResponseWriter, r *http.Request) {}
func UserStartFollowingDB(w http.ResponseWriter, r *http.Request)        {}
func UserStopFollowingDB(w http.ResponseWriter, r *http.Request)         {}

func TweetCreate(w http.ResponseWriter, r *http.Request) {}
func TweetRead(w http.ResponseWriter, r *http.Request)   {}
func TweetDelete(w http.ResponseWriter, r *http.Request) {}

//func TweetDeleteByCreatedTime(w http.ResponseWriter, r *http.Request) {}
