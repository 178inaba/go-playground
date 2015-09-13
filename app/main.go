package main

import (
	"flag"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"golang.org/x/oauth2"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

const settingToml = "setting/setting.toml"

var (
	s           setting
	mgoSessOrgn *mgo.Session
	oauth2Conf  = &oauth2.Config{
		Scopes: []string{"gist"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

type setting struct {
	server  `toml:"server"`
	mongo   `toml:"mongo"`
	client  `toml:"client"`
	sandbox `toml:"sandbox"`
}

type server struct {
	Port              int    `toml:"port"`
	OAuth2RedirectURL string `toml:"oauth2_redirect_url"`
}

type mongo struct {
	Host string `toml:"host"`
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
}

func main() {
	log.Debug("main()")

	loadSetting()

	// set oauth2 config
	oauth2Conf.ClientID = s.client.ID
	oauth2Conf.ClientSecret = s.client.Secret

	// connect mongo
	setMgoSess()

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
		log.Fatalf("setting load error: %v", err)
	}
}

func setMgoSess() {
	log.Debug("setMgoSess()")

	var err error
	mgoSessOrgn, err = mgo.Dial(s.mongo.Host)
	if err != nil {
		log.Fatalf("mongo dial error: %v", err)
	}

	mgoSessOrgn.SetMode(mgo.Monotonic, true)
}
