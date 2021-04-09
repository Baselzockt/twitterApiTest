package main

import (
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"twitterApiTest/twitterapi"
)

func setupLogging(loglevel log.Level, consoleOut bool) {
	log.SetLevel(loglevel)
	if consoleOut {
		log.SetOutput(os.Stdout)
		return
	}

	logfile, err0 := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err0 != nil {
		log.Fatal(err0)
	}

	log.SetOutput(logfile)
}

func main() {
	consoleLogging, _ := strconv.ParseBool(os.Getenv("CONSOLELOGGING"))
	setupLogging(log.DebugLevel, consoleLogging)
	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"Covid", "impfen"}}
	twitterapi.CreateHandlerForFilterStream(filterStreamParams)
	log.Debug("Received filter stream")
}
