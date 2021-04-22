package twitterapi

import (
	"bytes"
	"encoding/json"
	activeMq "github.com/Baselzockt/GoMQ/client"
	"github.com/Baselzockt/GoMQ/client/impl"
	"github.com/Baselzockt/GoMQ/content"
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
)

func CreateHandlerForFilterStream(twitterClient twitterClientInterface, activeMqClient activeMq.Client, filterParams *twitter.StreamFilterParams) error {
	stream, err := twitterClient.CreateFilterStream(filterParams)
	if err != nil {
		log.Error("Could not create filter stream")
		return err
	}
	log.Debug("Setting up twitter stream handler")
	err = setupTweetHandling(stream, activeMqClient)
	return err
}

func setupTweetHandling(stream *chan interface{}, client activeMq.Client) error {
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

	log.Debug("Starting handler")
	demux.HandleChan(*stream)
	return nil
}

func sendMessageToActiveMq(body []byte, client activeMq.Client) {
	log.Debug("sending message to activeMQ")

	err := client.SendMessageToQueue("Twitter", content.TEXT, body)
	if err != nil {
		log.Error(err)
		log.Debug("Error while sending trying to Reconnect")
		test := false
		switch client.(type) {
		case *impl.MockClient:
			err = client.Connect("localhost")
			test = true
		default:
			client = impl.NewStompClient()
			err = client.Connect("activemq:61613")
		}

		if err != nil {
			log.Error("Could not reconnect")
			if !test {
				log.Fatal(err)
			}
		}
	}
}
