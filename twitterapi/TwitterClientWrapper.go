package twitterapi

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
)

type twitterClient struct {
	client *twitter.Client
}

func NewTwitterClient(apiKey, apiSecret, accessKey, accessSecret string) *twitterClient {
	client := twitterClient{createClient(apiKey, apiSecret, accessKey, accessSecret)}
	return &client
}

func (t *twitterClient) CreateFilterStream(params *twitter.StreamFilterParams) (*chan interface{}, error) {
	stream, err := t.client.Streams.Filter(params)
	return &stream.Messages, err
}

func createClient(apiKey, apiSecret, accessKey, accessSecret string) *twitter.Client {
	log.Debug("Connecting to twitter api")
	config := oauth1.NewConfig(apiKey, apiSecret)
	token := oauth1.NewToken(accessKey, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	log.Debug("connected")
	return client
}
