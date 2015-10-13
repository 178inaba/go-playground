package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const (
	salt           = "The Go Playground to gist"
	maxSnippetSize = 64 * 1024
)

// Snippet is type of order to save the code made with playground.
type Snippet struct {
	ID   string
	Body []byte
}

func (snip *Snippet) setID() {
	h := sha1.New()
	io.WriteString(h, salt)
	h.Write(snip.Body)
	sum := h.Sum(nil)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(sum)))
	base64.URLEncoding.Encode(b, sum)

	snip.ID = string(b)[:10]
}

func init() {
	http.HandleFunc("/share", share)
}

func share(w http.ResponseWriter, r *http.Request) {
	id, err := saveSnip(w, r, "")
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Fprint(w, id)
}

func saveSnip(w http.ResponseWriter, r *http.Request, idSuffix string) (string, error) {
	if r.Method != "POST" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return "", errors.New("Forbidden")
	}

	var body bytes.Buffer
	_, err := io.Copy(&body, io.LimitReader(r.Body, maxSnippetSize+1))
	r.Body.Close()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return "", fmt.Errorf("reading Body: %v", err)
	}
	if body.Len() > maxSnippetSize {
		http.Error(w, "Snippet is too large", http.StatusRequestEntityTooLarge)
		return "", errors.New("Snippet is too large")
	}

	snip := &Snippet{Body: body.Bytes()}
	snip.setID()

	// add suffix
	if idSuffix != "" {
		snip.ID += "_" + idSuffix
	}

	// get mongo session
	mgoSess := mgoSessOrgn.Copy()
	defer mgoSess.Close()

	c := mgoSess.DB("playground").C("snippet")
	err = c.Insert(snip)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return "", fmt.Errorf("putting Snippet: %v", err)
	}

	return snip.ID, nil
}
