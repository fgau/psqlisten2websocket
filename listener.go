package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

const (
	host     = "your host"
	port     = 5432
	user     = "your db username"
	password = "your db password"
	dbname   = "your db name"
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
	var conninfo = fmt.Sprintf("host=%s dbname=%s user=%s password=%s", host, dbname, user, password)

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
