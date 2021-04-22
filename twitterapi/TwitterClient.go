package twitterapi

import "github.com/dghubble/go-twitter/twitter"

type TwitterClient interface {
	CreateFilterStream(*twitter.StreamFilterParams) (*chan interface{}, error)
}
