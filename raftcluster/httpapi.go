// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file is borrowed from https://github.com/etcd-io/etcd/tree/v3.3.20/contrib/raftexample
// I tweak it to adapt my requirement
package raftcluster

import (
    "log"
    "net/http"
    controller "zl2501-final-project/raftcluster/httpcontroller"
    "zl2501-final-project/raftcluster/mux"
)

func StartHttpAPI() {
    mux := mux.New()
    mux.Get("/hello", func(writer http.ResponseWriter, request *http.Request) {
        writer.Write([]byte("Hello! This is a New MUX"))
    })
    mux.Get("/auth/session/:sid", controller.GetSession)
    mux.Post("/auth/session/:sid",controller.CreateSession)
    log.Printf("Raft HTTP Server is going to start at: http://localhost:%v", HTTP_PORT)
    log.Fatal(http.ListenAndServe(":"+HTTP_PORT, mux))
}

// A middleware to log all requests
//func LogRequests(handlerToWrap http.Handler) http.Handler {
//    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//        logger := log.New(os.Stdout, "LogRequests:", log.Ltime|log.Lshortfile)
//        start := time.Now()
//        handlerToWrap.ServeHTTP(w, r)
//        logger.Printf("Request:%s %s, Time: %v", r.Method, r.URL.Path, time.Since(start))
//    })
//}

func MiddlewareAdapt(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
    for _, mw := range middleware {
        h = mw(h)
    }
    return h
}

