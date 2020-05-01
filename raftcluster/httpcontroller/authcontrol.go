package httpcontroller

import (
    "encoding/json"
    "net/http"
    "zl2501-final-project/raftcluster/mux"
    authstorage "zl2501-final-project/raftcluster/store/authstore"
    "zl2501-final-project/raftcluster/store/authstore/memory"
    _ "zl2501-final-project/raftcluster/store/authstore/memory"
)

var sessProvider authstorage.ProviderInterface

func init() {
    sessProvider, _ = authstorage.GetProvider("memory")
}
func GetSession(w http.ResponseWriter, r *http.Request) {
    sid := getParam(r, "sid")
    sess, err := sessProvider.SessionRead(sid)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    // The return should not be cached
    w.Header().Set("Cache-Control", "no-cache")
    w.WriteHeader(http.StatusCreated)
    sessIns := sess.(*memory.MemSessStore)
    ret, _ := json.Marshal(sessIns)
    w.Write(ret)
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
    sid := getParam(r, "sid")
    _, err := sessProvider.SessionInit(sid)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusNoContent)
}

// Params return the router params
func getParams(r *http.Request) map[string]string {
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
func getParam(r *http.Request, key string) string {
    p := getParams(r)
    if v, ok := p[key]; ok {
        return v
    }
    return ""
}
