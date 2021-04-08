package twitterapi

import (
	"bytes"
	"encoding/json"
	activeMq "github.com/Baselzockt/GoMQ/client"
	"github.com/Baselzockt/GoMQ/client/impl"
	"github.com/Baselzockt/GoMQ/content"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
)

var twitterClient *twitter.Client = nil

func getTwitterClient(securityConfig *SecurityConfig) *twitter.Client {

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

func CreateHandlerForFilterStream(filterParams *twitter.StreamFilterParams) {
	if twitterClient == nil {
		twitterClient = getTwitterClient(LoadSecurityConfig())
	}
	stream, err := twitterClient.Streams.Filter(filterParams)

	if err != nil {
		log.Error("Could not create Filter stream")
		log.Fatal(err)
	}

	log.Debug("Create and connect activeMQ client")
	activeMqClient := impl.NewStompClient()
	err = activeMqClient.Connect("activemq:61613")
	if err != nil {
		log.Error("Could not create create activemq client")
		log.Fatal(err)
	}

	log.Debug("Successfully created client")
	log.Debug("Setting up twitter stream handler")
	setupTweetHandling(stream, activeMqClient)
}

func setupTweetHandling(stream *twitter.Stream, client activeMq.Client) {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		err := encoder.Encode(tweet)

		if err != nil {
			log.Error("Could not encode tweet")
			log.Error(err)
			return
		}

		sendMessageToActiveMq(buffer.Bytes(), client)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		sendMessageToActiveMq([]byte(dm.Text), client)
	}
	log.Debug("Starting handler")
	demux.HandleChan(stream.Messages)
}

func sendMessageToActiveMq(body []byte, client activeMq.Client) {
	log.Debug("sending message to activeMQ")
	err := client.SendMessageToQueue("Twitter", content.TEXT, body)
	if err != nil {
		log.Error(err)
		log.Debug("Error while sending trying to Reconnect")
		client = impl.NewStompClient()
		err = client.Connect("activemq:61613")
		if err != nil {
			log.Fatal(err)
		}
	}
}
