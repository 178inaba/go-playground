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

	r.ParseForm()

	// get code
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

	// access token save in cookie
	log.Debugf("access token: %s", tok.AccessToken)

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: tok.AccessToken,
	})

	// get snip id
	_, ok = r.Form["snip_id"]
	if !ok {
		log.Error("snip id not found error")
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// snip id save in cookie
	log.Debugf("snip id: %s", r.Form["snip_id"][0])

	http.SetCookie(w, &http.Cookie{
		Name:  "snip_id",
		Value: r.Form["snip_id"][0],
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
