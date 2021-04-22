package main

import (
	"github.com/Baselzockt/GoMQ/client"
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
	consoleLogging, _ := strconv.ParseBool(os.Getenv("CONSOLELOGGING"))
	setupLogging(log.DebugLevel, consoleLogging)

	twitterClient := twitterapi.NewTwitterClient(os.Getenv("APIKEY"), os.Getenv("APISECRET"),
		os.Getenv("ACCESSKEY"), os.Getenv("ACCESSSECRET"))
	activeMqClient := impl.StompClient{Url: os.Getenv("ENDPOINT")}

	err := run(twitterClient, &activeMqClient)

	if err != nil {
		log.Fatal(err)
	}
}

func run(twitterClient twitterapi.TwitterClient, activeMqClient client.Client) error {
	log.Debug("Creating activeMq client")
	log.Debug("Connecting to activeMQ endpoint")
	err := activeMqClient.Connect()

	if err != nil {
		log.Error("Could not connect to activeMQ")
		return err
	}

	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"Covid", "impfen"}}

	return twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, filterStreamParams)
}
