package httpcontroller

import (
    "context"
    "encoding/json"
    "net/http"
    "zl2501-final-project/raftcluster/mux"
    . "zl2501-final-project/raftcluster/store"
    "zl2501-final-project/raftcluster/store/authstore/memory"
    _ "zl2501-final-project/raftcluster/store/authstore/memory"
)

func newTimeoutCtx() context.Context {
    ctx, _ := context.WithTimeout(context.Background(), ContextTimeoutDuration)
    return ctx
}
func ReadSession(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionRead,
        SessionProviderParams{sid})
    //sess, err := sessProvider.SessionRead(sid)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusCreated)
    sessIns := res.(*memory.MemSessStore)
    ret, _ := json.Marshal(requestRetType{
        "result": sessIns,
        "error":  nil,
    })
    w.Write(ret)
}
func CreateSession(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    _, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionInit,
        SessionProviderParams{sid})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    ret, _ := json.Marshal(requestRetType{
        "result": sid,
        "error": nil,
    })
    w.WriteHeader(http.StatusOK)
    w.Write(ret)
}
func SessionGC(w http.ResponseWriter, r *http.Request) {
    _, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionGC, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
func DestroySession(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionDestroy,
        SessionProviderParams{sid})
    //sess, err := sessProvider.SessionRead(sid)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusCreated)
    sessIns := res.(*memory.MemSessStore)
    ret, _ := json.Marshal(requestRetType{
        "result": sessIns,
        "error":  nil,
    })
    w.Write(ret)
}

func SessionGetValue(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    key := getRouteParam(r, "key")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionGet,
        SessionParams{Sid: sid, Key: key, Value: ""})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "error":  nil,
    })
    w.Write(ret)
}
func SessionPutKeyValue(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    key := getRouteParam(r, "key")
    r.ParseForm()
    value := r.Form["value"][0]
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionSet,
        SessionParams{Sid: sid, Key: key, Value: value})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "error":  nil,
    })
    w.Write(ret)
}
func SessionDeleteKey(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    key := getRouteParam(r, "key")
    res, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionDelete,
        SessionParams{Sid: sid, Key: key, Value: ""})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        ret, _ := json.Marshal(requestRetType{
            "result": nil,
            "error":  err.Error(),
        })
        w.Write(ret)
        return
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusOK)
    ret, _ := json.Marshal(requestRetType{
        "result": res,
        "error":  nil,
    })
    w.Write(ret)
}

// Params return the router params
func getRouteParams(r *http.Request) map[string]string {
    v := r.Context().Value(mux.ROUTE_PARAMS_KEY)
    if v == nil {
        return map[string]string{}
    }
    if v, ok := v.(map[string]string); ok {
        return v
    }
    return map[string]string{}
}

// Param return the router param based on the key
func getRouteParam(r *http.Request, key string) string {
    p := getRouteParams(r)
    if v, ok := p[key]; ok {
        return v
    }
    return ""
}
