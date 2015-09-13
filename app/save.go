package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func init() {
	http.HandleFunc("/save", saveHandler)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("saveHandler()")

	id, err := saveSnip(w, r, "save")
	if err != nil {
		log.Error(err)
		return
	}

	oauth2Conf.RedirectURL = s.OAuth2RedirectURL + "?snip_id=" + id
	fmt.Fprint(w, oauth2Conf.AuthCodeURL("state"))
}
