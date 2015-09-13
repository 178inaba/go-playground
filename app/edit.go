package main

import (
	"html/template"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
)

const hello = `package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")
}
`

var editTemplate = template.Must(template.ParseFiles("template/edit.html"))

type editData struct {
	Snippet *Snippet
}

func init() {
	http.HandleFunc("/", edit)
}

func edit(w http.ResponseWriter, r *http.Request) {
	// get mongo session
	mgoSess := mgoSessOrgn.Copy()
	defer mgoSess.Close()

	c := mgoSess.DB("playground").C("snippet")

	snip := &Snippet{Body: []byte(hello)}
	if strings.HasPrefix(r.URL.Path, "/p/") {
		id := r.URL.Path[3:]
		serveText := false
		if strings.HasSuffix(id, ".go") {
			id = id[:len(id)-3]
			serveText = true
		}

		err := c.Find(bson.M{"id": id}).One(&snip)
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
