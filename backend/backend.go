package backend

import (
	"log"
	//	"net/http"
	"zl2501-final-project/backend/constant"
)

func init() {
	// Set global logger
	log.SetPrefix("BE Service: ")
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func StartService() {
	//	mux := http.NewServeMux()
	log.Println("Server is going to start at: http://localhost:" + constant.Port)
}
