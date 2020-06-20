package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lib/pq"
)

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			log.Println("[dblistener] - received data from channel [", n.Channel, "] :")
			// prepare notification payload for pretty print
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				log.Println("[dblistener] - error processing json: ", err)
				return
			}
			fmt.Println(string(prettyJSON.Bytes()))
			h.broadcast <- string(prettyJSON.Bytes())
			return
		case <-time.After(90 * time.Second):
			log.Println("[dblistener] - received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}

func initDBListener() {
	type ConfigDatabase struct {
		Host     string `env:"HOST" env:"HOST"`
		Port     string `env:"PORT" env:"PORT"`
		User     string `env:"USER" env:"USER"`
		Password string `env:"PASSWORD" env:"PASSWORD"`
		DBName   string `env:"DBNAME" env:"DBNAME"`
	}

	var cfg ConfigDatabase

	errReadenv := cleanenv.ReadConfig(".env", &cfg)
	if errReadenv != nil {
		log.Println("[error] - cannot read env file", errReadenv)
	}

	var conninfo = fmt.Sprintf("host=%s dbname=%s user=%s password=%s", cfg.Host, cfg.DBName, cfg.User, cfg.Password)

	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("[dblistener] - " + err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	log.Println("[dblistener] - start monitoring postgresql ...")
	for {
		waitForNotification(listener)
	}
}
