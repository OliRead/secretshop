package main

import (
	"flag"
	"log"
	"net/http"

	"time"

	"github.com/oliread/secretshop"
	"github.com/oliread/secretshop/api"
	"github.com/oliread/secretshop/store/mysql"
)

func main() {
	confFile := flag.String("conf", "", "Location of config file")
	flag.Parse()

	conf, err := secretshop.ReadConfig(*confFile)
	if err != nil {
		log.Fatal(err)
	}

	if conf.Auth == "" {
		log.Print("WARNING: you are running Secret Shop without an authentication key set, be careful using this in the wild")
	}

	for host, data := range conf.StoreInfo {
		switch host {
		case "mysql":
			for i := 0; i < 5; i++ {
				if err := mysql.NewStore(&conf, data); err != nil {
					log.Printf("Attempt %d, Error connecting to database [%s]: %s", i+1, host, err)
				} else {
					break
				}
				time.Sleep(5 * time.Second)
			}
			if _, ok := conf.Stores[host]; !ok {
				log.Fatalf("Failed to connect to database [%s] after 5 attempts", host)
			}
		}
	}

	apiHandler, err := api.NewHandler(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving web endpoints on %s", conf.BindAddress)
	http.Handle("/", apiHandler.Router)

	log.Fatal(http.ListenAndServe(conf.BindAddress, http.DefaultServeMux))
}
