package main

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

func init() {
	http.HandleFunc("/gist", gistHandler)
}

func gistHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("gistHandler()")

	// get access token
	c, err := r.Cookie("access_token")
	if err != nil {
		log.Error("get cookie error: ", err)
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Value},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// github api
	client := github.NewClient(tc)

	// get post source
	byteBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("body read error: ", err)
	}

	body := string(byteBody)

	// create post gist
	postGist := github.Gist{
		Description: getStrPtr("test gist"),
		Public:      getBoolPtr(false),
		Files: map[github.GistFilename]github.GistFile{
			"test.go": github.GistFile{
				Content: &body,
			},
		},
	}
	retGist, resp, err := client.Gists.Create(&postGist)
	if err != nil {
		log.Error("create gist error: ", err)
		return
	}

	log.Debug(retGist)
	log.Debug(resp)
}

func getStrPtr(s string) *string {
	return &s
}

func getBoolPtr(b bool) *bool {
	return &b
}
