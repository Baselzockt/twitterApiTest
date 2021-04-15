package twitterapi

import "github.com/dghubble/go-twitter/twitter"

type twitterClientInterface interface {
	CreateFilterStream(*twitter.StreamFilterParams) (*chan interface{}, error)
}
