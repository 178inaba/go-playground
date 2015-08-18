package main

import (
	"html/template"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
)

func init() {
	http.HandleFunc("/", edit)
}

var editTemplate = template.Must(template.ParseFiles("edit.html"))

type editData struct {
	Snippet *Snippet
}

func edit(w http.ResponseWriter, r *http.Request) {
	// mongo
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("playground").C("snippet")

	snip := &Snippet{Body: []byte(hello)}
	if strings.HasPrefix(r.URL.Path, "/p/") {
		id := r.URL.Path[3:]
		serveText := false
		if strings.HasSuffix(id, ".go") {
			id = id[:len(id)-3]
			serveText = true
		}

		err = c.Find(bson.M{"id": id}).One(&snip)
		if err != nil {
			log.Errorf("loading Snippet: %v", err)
			http.Error(w, "Snippet not found", http.StatusNotFound)
			return
		}
		if serveText {
			w.Header().Set("Content-type", "text/plain")
			w.Write(snip.Body)
			return
		}
	}
	editTemplate.Execute(w, &editData{snip})
}

const hello = `package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")
}
`
