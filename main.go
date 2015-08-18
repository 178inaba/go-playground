package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const port = 8080

func main() {
	log.Info("main()")

	log.Info("server run: port: ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
