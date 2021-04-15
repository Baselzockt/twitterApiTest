package main

import (
	"github.com/Baselzockt/GoMQ/client/impl"
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"twitterApiTest/twitterapi"
)

func setupLogging(loglevel log.Level, consoleOut bool) {
	log.SetLevel(loglevel)
	log.SetReportCaller(true)
	if consoleOut {
		log.SetOutput(os.Stdout)
		return
	}

	var logfile, err = os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	log.RegisterExitHandler(func() {
		_ = logfile.Close()
	})
}

func main() {
	consoleLogging, _ := strconv.ParseBool(os.Getenv("CONSOLELOGGING"))
	setupLogging(log.DebugLevel, consoleLogging)
	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"Covid", "impfen"}}
	log.Debug("Creating twitter client")
	twitterClient := twitterapi.NewTwitterClient(os.Getenv("APIKEY"), os.Getenv("APISECRET"), os.Getenv("ACCESSKEY"), os.Getenv("ACCESSSECRET"))
	log.Debug("Creating activeMq client")
	activeMqClient := impl.NewStompClient()
	log.Debug("Connecting to activeMQ endpoint")
	err := activeMqClient.Connect(os.Getenv("ENDPOINT"))

	if err != nil {
		log.Error("Could not connect to activeMQ")
		log.Fatal(err)
	}

	err = twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, filterStreamParams)

	if err != nil {
		log.Error("Error while handling Filter stream")
		log.Fatal(err)
	}

	log.Debug("Received filter stream")
}
