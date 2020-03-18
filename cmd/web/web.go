package main

import (
	"log"
	"net/http"
	"zl2501-final-project/web"
)

func main() {
	println("Start the service")
	http.HandleFunc("/", web.SayhelloName)   // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
