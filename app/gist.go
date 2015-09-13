package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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

	gist, err := makeGist(r.Body)
	if err != nil {
		log.Errorf("make gist error: %v", err)
	}

	retGist, _, err := client.Gists.Create(&gist)
	if err != nil {
		log.Error("create gist error: ", err)
		return
	}

	io.WriteString(w, *retGist.HTMLURL)
}

func makeGist(httpBody io.Reader) (github.Gist, error) {
	// get body
	byteBody, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return github.Gist{}, fmt.Errorf("get body error: %v", err)
	}

	// unmarshal
	postGist := struct {
		FileName    string `json:"file_name"`
		Description string `json:"description"`
		Public      bool   `json:"public"`
		Code        string `json:"code"`
	}{}
	err = json.Unmarshal(byteBody, &postGist)
	if err != nil {
		return github.Gist{}, fmt.Errorf("json unmarshal error: %v", err)
	}

	// create gist
	gist := github.Gist{
		Description: github.String(strings.TrimSpace(postGist.Description)),
		Public:      &postGist.Public,
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(strings.TrimSpace(postGist.FileName) + ".go"): github.GistFile{
				Content: &postGist.Code,
			},
		},
	}

	return gist, nil
}
