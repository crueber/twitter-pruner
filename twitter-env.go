package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mkideal/cli"
)

// TwitterEnv is used to hold all twitter envs.
type TwitterEnv struct {
	cli.Helper
	ConsumerKey       string `cli:"*key" usage:"consumer key"`
	ConsumerSecret    string `cli:"*secret" usage:"consumer secret"`
	AccessToken       string `cli:"*token" usage:"access token"`
	AccessTokenSecret string `cli:"*tsecret" usage:"access token secret"`
	Days              int    `cli:"d,days" usage:"number of days to keep" dft:"28"`
	Rts               int    `cli:"rt" usage:"keep tweets with this many retweets" dft:"3"`
	Favs              int    `cli:"fav" usage:"keep tweets with this many favorites" dft:"3"`
	Commit            bool   `cli:"c" usage:"commit changes" dft:"false"`
}

// GenerateTwitterClient builds a twitter client that can be used to make calls
func (te *TwitterEnv) GenerateTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(te.ConsumerKey, te.ConsumerSecret)
	token := oauth1.NewToken(te.AccessToken, te.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}
