package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const port = 8080

func main() {
	log.Info("main()")

	serveStatic()

	log.Info("server run: port: ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func serveStatic() {
	dirName := "static"

	files, _ := ioutil.ReadDir(dirName)
	for _, f := range files {
		fileName := f.Name()
		http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, dirName+"/"+fileName)
		})
	}
}
