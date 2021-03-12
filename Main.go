package main

import (
	"bytes"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	content "twitterApiTest/activeMQ"
	activeMq "twitterApiTest/activeMQ/client"
	"twitterApiTest/activeMQ/client/impl"
	"twitterApiTest/twitterapi"
)

func setupLogging() {
	logfile, err0 := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err0 != nil {
		log.Fatal(err0)
	}
	log.SetLevel(log.DebugLevel)

	log.SetOutput(logfile)
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

		sendMessageToActiveMq(buffer.Bytes(), aClient)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		sendMessageToActiveMq([]byte(dm.Text), aClient)
	}
	log.Debug("Starting handler")
	go demux.HandleChan(stream.Messages)
}

func main() {
	setupLogging()
	client := getTwitterClient(loadSecurityConfig())

	log.Debug("getting filterstream")
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"ich","du","er","sie","es","Ich","Der","der","das","Das"}}
	stream, _ := client.Streams.Filter(filterStreamParams)
	log.Debug("Received filter stream")

	log.Debug("Create and connect activeMQ client")
	aClient := impl.NewStompClient()
	err := aClient.Connect("localhost:61613")
	if err == nil {
		log.Debug("Successfully created client")
		log.Debug("Setting up twitter stream handler")
		setupTweetHandling(stream, aClient)

		log.Debug("Subscribing to queue")
		aClient.SubscribeToQueue("Twitter", channel)

		log.Debug("setting up websocket endpoint")
		http.HandleFunc("/", wsEndpoint)
		log.Fatal(http.ListenAndServe(":5000", nil))
	}
}

var channel = make(chan []byte)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func sendMessageToActiveMq(body []byte, client activeMq.Client) {
	log.Debug("sending message to activeMQ")
	client.SendMessageToQueue("Twitter", content.TEXT, body)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	log.Debug(" trying to upgrade to ws")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("Upgrade unsuccessful")
	} else {
		log.Debug("Client connected")
		for {
			if !handleClient(ws) {
				break
			}
		}
		log.Debug("close connection")
		ws.Close()
	}
}

func handleClient(ws *websocket.Conn) bool{
	_, response, err := ws.ReadMessage()
	log.Debug("got response: ", string(response))
	if err != nil || string(response) != "ok" {
		log.Debug("closing connection")
		ws.WriteMessage(1, []byte("Closing connection"))
		return false
	}
	log.Debug("sending answer")
	err2 := ws.WriteMessage(1, <-channel)
	if err2 != nil {
		log.Debug("could not send answer because of error: " + err2.Error())
		return false
	}
	return true
}
