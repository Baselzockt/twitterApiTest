package twitterapi

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
)

type TwitterClientMock struct {
	tweets []twitter.Tweet
}

func NewTwitterClientMock(tweets []twitter.Tweet) *TwitterClientMock {
	return &TwitterClientMock{tweets: tweets}
}

func (t *TwitterClientMock) CreateFilterStream(params *twitter.StreamFilterParams) (*chan interface{}, error) {
	if t.tweets == nil {
		return nil, errors.New("could not initiate Filter Stream")
	}

	msgChan := make(chan interface{})
	go func(msgChan chan interface{}) {
		for _, tweet := range t.tweets {
			msgChan <- &tweet
		}
		close(msgChan)
	}(msgChan)
	return &msgChan, nil
}
