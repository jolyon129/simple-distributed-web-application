package controller

import (
    "context"
    "golang.org/x/crypto/bcrypt"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "time"
    "zl2501-final-project/web/constant"
    . "zl2501-final-project/web/pb"
)

func EncodePassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    } else {
        return string(hash), nil
    }
}

func ComparePassword(p1 string, p2 string) error {
    return bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
}

// Read sessionId from cookie If existed.
// If not exist, create a new sessionId and inject sessId into cookie.
// If exist and the sessionId is valid, reuse the same session.
func SessionStart(w http.ResponseWriter, r *http.Request) (string, error) {
    cookie, err := r.Cookie(constant.SessCookieName)
    ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
    var newSessId string
    if err != nil || cookie.Value == "" {
        response, err := AuthClientIns.SessionStart(ctx1, &SessionGeneralRequest{
            SessionId: "", // new session with empty sessId
        })
        if err != nil {
            return "", err
        }
        newSessId = response.SessionId
    } else {
        sessId, _ := url.QueryUnescape(cookie.Value)
        response, err := AuthClientIns.SessionStart(ctx1, &SessionGeneralRequest{
            SessionId: sessId, // try to reuse session
        })
        if err != nil {
            log.Print(err)
        }
        newSessId = response.SessionId // Use the new one
    }
    newCookie := http.Cookie{Name: constant.SessCookieName, Value: url.QueryEscape(newSessId), Path: "/",
        HttpOnly: true, MaxAge: int(constant.MaxLifeTime)}
    http.SetCookie(w, &newCookie)
    return newSessId, nil
}

// Initialize a new session regardless of the cookie
func SessionInit(w http.ResponseWriter, r *http.Request) (string,error){
    ctx1, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
    response, err := AuthClientIns.SessionStart(ctx1, &SessionGeneralRequest{
        SessionId: "", // set a new session
    })
    if err != nil {
        log.Print(err)
    }
    newSessId := response.SessionId // Use the new one
    newCookie := http.Cookie{Name: constant.SessCookieName, Value: url.QueryEscape(newSessId), Path: "/",
        HttpOnly: true, MaxAge: int(constant.MaxLifeTime)}
    http.SetCookie(w, &newCookie)
    return newSessId, nil

}

func GetMyUserId(r *http.Request) (uint64, error) {
    cookie, _ := r.Cookie(constant.SessCookieName)
    sessId, _ := url.QueryUnescape(cookie.Value)
    ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
    getValueResponse, err := AuthClientIns.GetValue(ctx, &GetValueRequest{
        Key:  constant.UserId,
        Ssid: sessId,
    })
    if err != nil {
        return 0, err
    }
    userId := getValueResponse.Value
    myUidint, _ := strconv.Atoi(userId)
    myUid := uint64(myUidint)
    return myUid, nil
}

// Auth the request
func CheckAuthRequest(r *http.Request) bool {
    ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
    cookie, err := r.Cookie(constant.SessCookieName)
    if err != nil || cookie.Value == "" {
        log.Printf("Request:%s %s is not authenticated. Redirect to index.", r.Method, r.URL.Path)
        return false
    } else {
        sid, _ := url.QueryUnescape(cookie.Value)
        response, err := AuthClientIns.SessionAuth(ctx, &SessionGeneralRequest{
            SessionId: sid,
        })
        if err != nil {
            log.Printf(err.Error()) //TODO: error Handler
            return false
        }
        if response.Ok {
            log.Printf("Request:%s %s is authenticated.", r.Method,
                r.URL.Path)
            return true
        }
        return false
    }
}

// Manually terminate the session and overwrite the corresponding cookie into empty
func SessionDestroy(w http.ResponseWriter, r *http.Request) {
    ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
    cookie, err := r.Cookie(constant.SessCookieName)
    if err != nil || cookie.Value == "" {
        log.Printf("Request:%s %s is not authenticated.", r.Method, r.URL.Path)
    } else {
        sid, _ := url.QueryUnescape(cookie.Value)
        response, err := AuthClientIns.SessionDestroy(ctx, &SessionGeneralRequest{
            SessionId: sid,
        })
        if err != nil {
            log.Printf(err.Error()) //TODO: error Handler
        }
        if response.Ok {
            log.Printf("Request:%s %s is complete.", r.Method,
                r.URL.Path)
        }
    }
    expiration := time.Now()
    // Set empty cookie
    cookie1 := http.Cookie{Name: constant.SessCookieName, Path: "/", HttpOnly: true,
        Expires: expiration, MaxAge: -1}
    http.SetCookie(w, &cookie1)

}
