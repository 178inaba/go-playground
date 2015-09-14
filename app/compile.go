package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

func init() {
	http.HandleFunc("/compile", compile)
}

func compile(w http.ResponseWriter, r *http.Request) {
	if err := passThru(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Compile server error.")
	}
}

func passThru(w io.Writer, req *http.Request) error {
	log.Debug("passThru()")

	defer req.Body.Close()

	jsonReader, err := makeBodyJSON(req.Body)
	if err != nil {
		log.Errorf("make body json error: %v", err)
		return err
	}

	r, err := http.Post(s.sandbox.URL, req.Header.Get("Content-type"), jsonReader)
	if err != nil {
		log.Errorf("making POST request: %v", err)
		return err
	}
	defer r.Body.Close()
	if _, err := io.Copy(w, r.Body); err != nil {
		log.Errorf("copying response Body: %v", err)
		return err
	}
	return nil
}

func makeBodyJSON(httpBody io.Reader) (io.Reader, error) {
	// io.Reader -> []byte
	httpBodyByte, err := ioutil.ReadAll(httpBody)
	if err != nil {
		log.Errorf("body read error: %v", err)
		return nil, err
	}

	// []byte -> url.Values
	v, err := url.ParseQuery(string(httpBodyByte))
	if err != nil {
		log.Errorf("query parse error: %v", err)
		return nil, err
	}

	// make json
	json, err := json.Marshal(map[string]string{"Body": v.Get("body")})
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil, err
	}

	// convert byte json -> io.Reader
	return bytes.NewReader(json), nil
}
