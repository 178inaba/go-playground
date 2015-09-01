package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func init() {
	http.HandleFunc("/gist", gistHandler)
}

func gistHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("gistHandler()")

	c, err := r.Cookie("access_token")
	if err != nil {
		log.Error("get cookie error: ", err)
	}

	log.Debug(c.Value)

	// TODO(178inaba): create gist
}
