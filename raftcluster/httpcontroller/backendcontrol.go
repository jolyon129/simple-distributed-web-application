package httpcontroller

import (
    "encoding/json"
    "net/http"
    "strconv"
    . "zl2501-final-project/raftcluster/store"
    "zl2501-final-project/raftcluster/store/backendstore"
)

func convertStrToUint(s string) uint {
    myUidint, _ := strconv.Atoi(s)
    return uint(myUidint)
}

type requestRetType map[string]interface{}

// Get the user info from request.from
func UserCreate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    userName := r.Form["username"][0]
    password := r.Form["password"][0]
    uid, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserCreate,
        UserCreateParams{UserInfo{
            UserName: userName,
            Password: password,
        }})
    var ret []byte
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ = json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
    } else {
        w.WriteHeader(http.StatusCreated)
        ret, _ = json.Marshal(requestRetType{
            "result": uid,
            "error":  err,
        })
    }
    w.Write(ret)
}

//func UserDelete(w http.ResponseWriter, r *http.Request){}

func UserRead(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    uid := getRouteParam(r, "uid")
    nUid := convertStrToUint(uid)
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserGet,
        UserIDParams{nUid})
    w.Header().Set("Cache-Control", "no-cache")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.WriteHeader(http.StatusOK)
    user := res.(*backendstore.UserEntity)
    ret, _ := json.Marshal(user)
    w.Write(ret)
}

//func UserUpdate(w http.ResponseWriter, r *http.Request) {
//}

func UserFindAll(w http.ResponseWriter, r *http.Request) {
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserFindAll, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    users := res.([]*backendstore.UserEntity)
    ret, _ := json.Marshal(requestRetType{
        "users": users,
        "error": err,
    })
    w.Write(ret)
}

func UserAddTweetToUserDB(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    uid := getRouteParam(r, "uid")
    tid := getRouteParam(r, "tid")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserAddTweetToUserDB,
        UserAddTweetToUserParams{
            UId: convertStrToUint(uid),
            TId: convertStrToUint(tid),
        })
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "err":    err,
    })
    w.Write(ret)
}
func UserCheckWhetherFollowingDB(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    srcId := getRouteParam(r, "srcuid")
    tarId := getRouteParam(r, "targetuid")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserCheckWhetherFollowingGetDB,
        UserCheckWhetherFollowingDBParams{
            SrcId:    convertStrToUint(srcId),
            TargetId: convertStrToUint(tarId),
        })
    w.Header().Set("Cache-Control", "no-cache")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "error":  err,
    })
    w.Write(ret)
}
func UserStartFollowingDB(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    srcId := getRouteParam(r, "srcuid")
    tarId := getRouteParam(r, "targetuid")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserStartFollowingDB,
        UserCheckWhetherFollowingDBParams{
            SrcId:    convertStrToUint(srcId),
            TargetId: convertStrToUint(tarId),
        })
    w.Header().Set("Cache-Control", "no-cache")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "error":  err,
    })
    w.Write(ret)
}
func UserStopFollowingDB(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    srcId := getRouteParam(r, "srcuid")
    tarId := getRouteParam(r, "targetuid")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_UserStopFollowingDB,
        UserCheckWhetherFollowingDBParams{
            SrcId:    convertStrToUint(srcId),
            TargetId: convertStrToUint(tarId),
        })
    w.Header().Set("Cache-Control", "no-cache")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "error":  err,
    })
    w.Write(ret)

}

func TweetCreate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    uId := r.Form["uid"][0]
    content := r.Form["content"][0]
    tid, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_TweetCreate, TweetInfo{
        UserID:  convertStrToUint(uId),
        Content: content,
    })
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }

    w.WriteHeader(http.StatusCreated)
    nres := requestRetType{
        "result": tid,
        "error":  nil,
    }
    ret, _ := json.Marshal(nres)
    w.Write(ret)
}

func TweetRead(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    tid := getRouteParam(r, "tid")
    nTid := convertStrToUint(tid)
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_TweetGet, TweetReadParams{nTid})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }

    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    tweet := res.(*backendstore.TweetEntity)
    ret, _ := json.Marshal(requestRetType{
        "result": tweet,
        "error":  err,
    })
    w.Write(ret)
}

func TweetGetAll(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_TweetGetAll, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err,
        })
        w.Write(ret)
        return
    }
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    tweet := res.([]*backendstore.TweetEntity)
    ret, _ := json.Marshal(requestRetType{
        "result": tweet,
        "error":  err,
    })
    w.Write(ret)
}

//func TweetDelete(w http.ResponseWriter, r *http.Request) {}

//func TweetDeleteByCreatedTime(w http.ResponseWriter, r *http.Request) {}
