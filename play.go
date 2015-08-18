package main

import (
	"io"
	"net/http"

	"golang.org/x/tools/godoc/static"
)

func init() {
	http.HandleFunc("/playground.js", play)
}

func play(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/javascript")
	io.WriteString(w, static.Files["playground.js"])
}
