package tests

import (
	"github.com/Baselzockt/GoMQ/client/impl"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
	"testing"
	"twitterApiTest/twitterapi"
)

func TestTweetHandling(t *testing.T) {
	filterStreamParams := &twitter.StreamFilterParams{Language: []string{"de"}, Track: []string{"Covid", "impfen"}}

	t.Run("Happy path", func(t *testing.T) {
		twitterClient := twitterapi.NewTwitterClientMock(createTestData(50))
		activeMqClient := impl.NewMockClient()
		err := activeMqClient.Connect("localhost")
		assert.Nil(t, err)
		err = twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, filterStreamParams)
		assert.Nil(t, err)
		assert.True(t, len(activeMqClient.GetMessages()) == 50)
	})

	t.Run("No tweets", func(t *testing.T) {
		twitterClient := twitterapi.NewTwitterClientMock(createTestData(0))
		activeMqClient := impl.NewMockClient()
		err := activeMqClient.Connect("localhost")
		assert.Nil(t, err)
		err = twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, filterStreamParams)
		assert.Nil(t, err)
		assert.True(t, len(activeMqClient.GetMessages()) == 0)
	})

	t.Run("test nil", func(t *testing.T) {
		twitterClient := twitterapi.NewTwitterClientMock(nil)
		activeMqClient := impl.NewMockClient()
		err := twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, nil)
		assert.NotNil(t, err)
		assert.True(t, len(activeMqClient.GetMessages()) == 0)
	})

	t.Run("test not connected", func(t *testing.T) {
		twitterClient := twitterapi.NewTwitterClientMock(createTestData(50))
		activeMqClient := impl.NewMockClient()
		err := twitterapi.CreateHandlerForFilterStream(twitterClient, activeMqClient, filterStreamParams)
		assert.Nil(t, err)
	})

}

func createTestData(tweetAmount int) []twitter.Tweet {
	var data = make([]twitter.Tweet, tweetAmount)
	for i := 0; i < tweetAmount; i++ {
		data[i] = twitter.Tweet{Lang: "de", Text: "Covid ist doof, aber wir können zum glück ja impfen"}
	}
	return data
}
