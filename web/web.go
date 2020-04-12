package web

import (
    "context"
    "google.golang.org/grpc"
    "log"
    "net/http"
    "time"
    "zl2501-final-project/web/constant"
    "zl2501-final-project/web/pb"
)

// Then, initialize the session manager
func init() {
    // Set global logger
    log.SetPrefix("GlobalLogger: ")
    log.SetFlags(log.Ltime | log.Lshortfile)
    // Set up a connection to the server.

}

func StartService() {
    conn, err := grpc.Dial(constant.BackendServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    backendClient := pb.NewBackendClient(conn)

    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := backendClient.SayHello(ctx, &pb.HelloRequest{Name: "Zhuolun"})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.GetMessage())
    //mux := http.NewServeMux()
    //mux.Handle("/", MiddlewareAdapt(http.HandlerFunc(controller.GoIndex), SetHeader))                                  // set router
    //mux.Handle("/index", MiddlewareAdapt(http.HandlerFunc(controller.GoIndex), SetHeader))
    //mux.Handle("/login", MiddlewareAdapt(http.HandlerFunc(controller.LogIn), SetHeader))
    //mux.Handle("/signup", MiddlewareAdapt(http.HandlerFunc(controller.SignUp), SetHeader))
    //mux.Handle("/home", MiddlewareAdapt(http.HandlerFunc(controller.Home), auth.CheckAuth, SetHeader))
    //mux.Handle("/logout", MiddlewareAdapt(http.HandlerFunc(controller.LogOut), SetHeader))
    //mux.Handle("/tweet", MiddlewareAdapt(http.HandlerFunc(controller.Tweet), auth.CheckAuth, SetHeader))
    //mux.Handle("/users", MiddlewareAdapt(http.HandlerFunc(controller.ViewUsers), auth.CheckAuth, SetHeader))
    //mux.Handle("/user/", MiddlewareAdapt(http.HandlerFunc(controller.User), auth.CheckAuth, SetHeader))
    //mux.Handle("/follow", MiddlewareAdapt(http.HandlerFunc(controller.Follow), auth.CheckAuth, SetHeader))
    //mux.Handle("/unfollow", MiddlewareAdapt(http.HandlerFunc(controller.Unfollow), auth.CheckAuth, SetHeader))
    //log.Println("Server is going to start at: http://localhost:"+constant.Port)
    //log.Fatal(http.ListenAndServe(":"+constant.Port, logger.LogRequests(mux)))
}

// Adapt all middleware to the handler.
// The function will call them one by one (in reverse order) in a chained manner,
// returning the result of the first adapter.
// Ref: https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
func MiddlewareAdapt(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
    for _, mw := range middleware {
        h = mw(h)
    }
    return h
}

// This is a middleware to
// add Some Header to response
func SetHeader(handlerToWrap http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.Header().Set("cache-control", "no-store")
        handlerToWrap.ServeHTTP(w, r)
    })
}
