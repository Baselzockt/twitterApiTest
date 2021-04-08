package main

import (
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
	"os"
	"twitterApiTest/twitterapi"
)

func setupLogging(loglevel log.Level) {
	logfile, err0 := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err0 != nil {
		log.Fatal(err0)
	}
	log.SetLevel(loglevel)
	log.SetOutput(logfile)
}

func main() {
	setupLogging(log.DebugLevel)
	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"Covid", "impfen"}}
	twitterapi.CreateHandlerForFilterStream(filterStreamParams)
	log.Debug("Received filter stream")
}
