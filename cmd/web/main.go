package main

import (
	"os"
	"zl2501-final-project/web"
)

func main() {
	// Change working directory to web/ first
	err := os.Chdir("../../web")
	if err != nil {
		panic(err)
	}
	web.StartService()
}
