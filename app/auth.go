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
	_, ok := r.Form["code"]
	if !ok {
		log.Error("auth code not found error")
		http.Error(w, "Github Auth Error", http.StatusInternalServerError)
		return
	}

	tok, err := oauth2Conf.Exchange(oauth2.NoContext, r.Form["code"][0])
	if err != nil {
		log.Errorf("code exchange error: %v", err)
		http.Error(w, "Github Auth Error", http.StatusInternalServerError)
		return
	}

	// access token
	log.Debugf("access token: %s", tok.AccessToken)

	c := http.Cookie{
		Name:  "access_token",
		Value: tok.AccessToken,
	}
	http.SetCookie(w, &c)

	http.Redirect(w, r, "/", http.StatusFound)
}
