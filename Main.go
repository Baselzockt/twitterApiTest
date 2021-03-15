package main

import (
	"bytes"
	"encoding/json"
	activeMq "github.com/Baselzockt/GoMQ/client"
	"github.com/Baselzockt/GoMQ/client/impl"
	"github.com/Baselzockt/GoMQ/content"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"os"
	"twitterApiTest/twitterapi"
)

func setupLogging(loglevel log.Level) {
	//logfile, err0 := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	/*if err0 != nil {
		log.Fatal(err0)
	}*/
	log.SetLevel(loglevel)

	log.SetOutput(os.Stdout)
}

func loadSecurityConfig() *twitterapi.SecurityConfig {
	log.Debug("Loading SecurityConfig from conf.json")
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	securityConfig := &twitterapi.SecurityConfig{}
	err1 := decoder.Decode(&securityConfig)

	if err1 != nil {
		log.Error("Could not load config creating file and closing program")
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.Encode(securityConfig)
		os.WriteFile("conf.json", buffer.Bytes(), os.ModeExclusive)
		log.Fatal("Config file could not be parsed")
	}
	file.Close()

	return securityConfig
}

func getTwitterClient(securityConfig *twitterapi.SecurityConfig) *twitter.Client {

	if securityConfig.AccessSecret == "" {
		log.Fatal("Config is empty")
	}

	log.Debug("Connecting to twitter api")
	config := oauth1.NewConfig(securityConfig.ApiKey, securityConfig.ApiSecret)
	token := oauth1.NewToken(securityConfig.AccessToken, securityConfig.AccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	log.Debug("connected")
	return client
}

func setupTweetHandling(stream *twitter.Stream, aClient activeMq.Client) {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.Encode(tweet)

		sendMessageToActiveMq(buffer.Bytes())
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		sendMessageToActiveMq([]byte(dm.Text))
	}
	log.Debug("Starting handler")
	demux.HandleChan(stream.Messages)
}

func main() {
	setupLogging(log.DebugLevel)
	client := getTwitterClient(loadSecurityConfig())

	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"ich", "du", "er", "sie", "es", "Ich", "Der", "der", "das", "Das", "Covid", "impfen"}}
	stream, _ := client.Streams.Filter(filterStreamParams)
	log.Debug("Received filter stream")

	log.Debug("Create and connect activeMQ client")
	aClient := impl.NewStompClient()
	err := aClient.Connect("activemq:61613")
	if err == nil {
		log.Debug("Successfully created client")
		log.Debug("Setting up twitter stream handler")
		setupTweetHandling(stream, aClient)
	}
}

var aClient = impl.NewStompClient()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func sendMessageToActiveMq(body []byte) {
	log.Debug("sending message to activeMQ")
	err := aClient.SendMessageToQueue("Twitter", content.TEXT, body)
	if err != nil {
		log.Debug("Reconnecting...")
		aClient = impl.NewStompClient()
		aClient.Connect("activemq:61613")
		log.Error(err)
	}
}
