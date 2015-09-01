package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

const (
	port        = 8080
	settingToml = "setting.toml"
)

var (
	s          setting
	oauth2Conf = &oauth2.Config{
		Scopes: []string{"gist"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

type setting struct {
	client `toml:"client"`
}

type client struct {
	Id     string `toml:"id"`
	Secret string `toml:"secret"`
}

func init() {
	log.SetLevel(log.DebugLevel)

	loadSetting()

	oauth2Conf.ClientID = s.client.Id
	oauth2Conf.ClientSecret = s.client.Secret
}

func main() {
	log.Info("main()")

	serveStatic()

	log.Info("server run: port: ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func serveStatic() {
	dirName := "static"

	files, _ := ioutil.ReadDir(dirName)
	for _, f := range files {
		fileName := f.Name()
		http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, dirName+"/"+fileName)
		})
	}
}

func loadSetting() {
	log.Debug("loadSetting()")

	_, err := toml.DecodeFile(settingToml, &s)
	if err != nil {
		log.Error("setting load error: ", err)
	}
}
