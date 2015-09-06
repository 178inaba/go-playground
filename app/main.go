package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

const settingToml = "setting/setting.toml"

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
	server  `toml:"server"`
	client  `toml:"client"`
	sandbox `toml:"sandbox"`
}

type server struct {
	Port int `toml:"port"`
}

type client struct {
	ID     string `toml:"id"`
	Secret string `toml:"secret"`
}

type sandbox struct {
	URL string `toml:"url"`
}

func init() {
	debug := flag.Bool("d", false, "output debug log.")
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	loadSetting()

	oauth2Conf.ClientID = s.client.ID
	oauth2Conf.ClientSecret = s.client.Secret
}

func main() {
	log.Info("main()")

	http.Handle("/js/", http.FileServer(http.Dir("static")))
	http.Handle("/css/", http.FileServer(http.Dir("static")))
	http.Handle("/img/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/favicon.ico")
	})

	log.Info("server run: port: ", s.server.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.server.Port), nil)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loadSetting() {
	log.Debug("loadSetting()")

	_, err := toml.DecodeFile(settingToml, &s)
	if err != nil {
		log.Error("setting load error: ", err)
	}
}
