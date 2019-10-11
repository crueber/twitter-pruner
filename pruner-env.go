package main

import (
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mkideal/cli"
)

// PrunerEnv is used to hold all twitter envs.
type PrunerEnv struct {
	cli.Helper
	ConsumerKey         string    `cli:"*key" usage:"consumer key"`
	ConsumerSecret      string    `cli:"*secret" usage:"consumer secret"`
	AccessToken         string    `cli:"*token" usage:"access token"`
	AccessTokenSecret   string    `cli:"*tsecret" usage:"access token secret"`
	Days                int       `cli:"d,days" usage:"number of days to keep" dft:"28"`
	Rts                 int       `cli:"rt" usage:"keep tweets with this many retweets" dft:"3"`
	Favs                int       `cli:"fav" usage:"keep tweets with this many favorites" dft:"3"`
	AllRts              bool      `cli:"r,allrts" usage:"removal all of your retweets" dft:"true"`
	Commit              bool      `cli:"c" usage:"commit changes" dft:"false"`
	MaxAPITweets        int       `cli:"max" usage:"max api tweets" dft:"3200"`
	MaxTweetsPerRequest int       `cli:"request" usage:"number of tweets per request" dft:"100"`
	MaxAge              time.Time `cli:"age" usage:"specific date that overrides days duration"`
}

// GenerateTwitterClient builds a twitter client that can be used to make calls
func (te *PrunerEnv) GenerateTwitterClient() (*twitter.Client, error) {
	if te.MaxAge.IsZero() {
		age, err := time.ParseDuration("-" + strconv.Itoa(te.Days*24) + "h")
		if err != nil {
			return nil, err
		}
		te.MaxAge = time.Now().Add(age)
	}

	config := oauth1.NewConfig(te.ConsumerKey, te.ConsumerSecret)
	token := oauth1.NewToken(te.AccessToken, te.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient), nil
}
