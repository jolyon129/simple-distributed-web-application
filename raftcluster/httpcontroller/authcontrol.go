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
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusCreated)
    sessIns := res.(*memory.MemSessStore)
    ret, _ := json.Marshal(sessIns)
    w.Write(ret)
}
func CreateSession(w http.ResponseWriter, r *http.Request) {
    sid := getRouteParam(r, "sid")
    _, err := raftStore.RequestPropose(newTimeoutCtx(), METHOD_SessionInit,
        SessionProviderParams{sid})
    //_, err := sessProvider.SessionInit(sid)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusNoContent)
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
