package auth

import (
    "context"
    "log"
    "net/http"
    "net/url"
    "zl2501-final-project/web"
    "zl2501-final-project/web/constant"
    "zl2501-final-project/web/pb"
)

// This is a middleware handler used to check weather this request is authenticated.
// If not, redirect to the index.
func CheckAuth(handlerToWrap http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //web.AuthClient
        ctx, _ := context.WithTimeout(context.Background(), constant.ContextTimeoutDuration)
        cookie, err := r.Cookie(constant.CookieName)
        if err != nil || cookie.Value == "" {
            log.Printf("Request:%s %s is not authenticated. Redirect to index.", r.Method, r.URL.Path)
            http.Redirect(w, r, "/", 302) // Go the index
        } else {
            sid, _ := url.QueryUnescape(cookie.Value)
            response, err := web.AuthClient.SessionAuth(ctx, &pb.SessionGeneralRequest{
                SessionId: sid,
            })
            if err != nil {
                log.Printf("Request:%s %s is not authenticated. Redirect to index.", r.Method, r.URL.Path)
                http.Redirect(w, r, "/", 302) // Go the index
                return
            }

            handlerToWrap.ServeHTTP(w, r)

        }

        //if ok {
        //	handlerToWrap.ServeHTTP(w, r)
        //} else {
        //	log.Printf("Request:%s %s is not authenticated. Redirect to index.", r.Method, r.URL.Path)
        //	http.Redirect(w, r, "/", 302) // Go the index
        //}
    })
}
