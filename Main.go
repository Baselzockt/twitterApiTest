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

	var logfile, err = os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	log.RegisterExitHandler(func() {
		logfile.Close()
	})
}

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	consoleLogging, _ := strconv.ParseBool(os.Getenv("CONSOLELOGGING"))
	setupLogging(log.DebugLevel, consoleLogging)

	log.Debug("Creating twitter client")
	twitterClient := twitterapi.NewTwitterClient(os.Getenv("APIKEY"), os.Getenv("APISECRET"), os.Getenv("ACCESSKEY"), os.Getenv("ACCESSSECRET"))

	log.Debug("Creating activeMq client")
	impl.StompClient{Url: os.Getenv("ENDPOINT")}
	activeMqClient := impl.NewStompClient()
	log.Debug("Connecting to activeMQ endpoint")
	err := activeMqClient.Connect(os.Getenv("ENDPOINT"))

	if err != nil {
		log.Error("Could not connect to activeMQ")
		return err
	}

	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"Covid", "impfen"}}

	return twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, filterStreamParams)
}
