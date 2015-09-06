package main

import (
	"net/http"

	"golang.org/x/oauth2"

	log "github.com/Sirupsen/logrus"
)

func init() {
	http.HandleFunc("/auth", authHandler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("authHandler()")

	// get code
	r.ParseForm()
	tok, err := oauth2Conf.Exchange(oauth2.NoContext, r.Form["code"][0])
	if err != nil {
		log.Fatal(err)
	}

	// access token
	log.Info("access token: ", tok.AccessToken)

	c := http.Cookie{
		Name:  "access_token",
		Value: tok.AccessToken,
	}
	http.SetCookie(w, &c)

	http.Redirect(w, r, "/", 303)
}
