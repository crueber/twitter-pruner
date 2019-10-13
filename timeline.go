package main

import (
	"fmt"
	"time"

	"github.com/crueber/twitter-pruner/pruner"
	"github.com/dghubble/go-twitter/twitter"
)

func isAgedOut(t *twitter.Tweet, env *pruner.Env) bool {
	createdTime, _ := t.CreatedAtTime()

	return env.MaxAge.After(createdTime)
}

func isBoring(t *twitter.Tweet, env *pruner.Env) bool {
	if env.AllRts && t.Retweeted {
		return true
	}
	if t.FavoriteCount >= env.Favs || t.RetweetCount >= env.Rts {
		rt := ""
		if t.Retweeted {
			rt += "re"
		}
		if env.Verbose {
			fmt.Printf("Ignoring %vtweet (%v fav/%v rt): %v\n", rt, t.FavoriteCount, t.RetweetCount, t.Text)
		}
		return false
	}
	return true
}

func calcTweetsToDelete(tweets []twitter.Tweet, env *pruner.Env) []int64 {
	var tweetsToDelete []int64
	for _, tweet := range tweets {
		if isAgedOut(&tweet, env) && isBoring(&tweet, env) {
			if env.Verbose {
				fmt.Printf("%v --- %v %v --- %v\n", tweet.CreatedAt, tweet.FavoriteCount, tweet.RetweetCount, tweet.Text)
			}
			tweetsToDelete = append(tweetsToDelete, tweet.ID)
		}
	}
	return tweetsToDelete
}

func deleteTweets(c *pruner.Client, tweetIds []int64) (int, int) {
	count := 0
	errorCount := 0
	if c.Env.Commit {
		for _, id := range tweetIds {
			err := c.DestroyTweet(id)
			if err != nil {
				if c.Env.Verbose {
					fmt.Printf("\n")
				}
				fmt.Printf("Error removing status: %v\n", err)
				errorCount++
				continue
			}
			if c.Env.Verbose {
				fmt.Printf(".")
			}
			count++
		}
	}
	return count, errorCount
}

// PruneTimeline does exactly what it says it does
func PruneTimeline(c *pruner.Client, user *twitter.User) error {
	var max int64
	count := 0
	markedForRemoval := 0
	removed := 0
	errorCount := 0
	shouldContinue := true

	for shouldContinue {
		c.Env.MaxAPICalls--
		tweets, err := c.GetTimeline(max)
		if err != nil {
			fmt.Printf("Error in timeline retrieval: %+v", err)
			errorCount++
		}
		count += len(tweets)

		tweetsToDelete := calcTweetsToDelete(tweets, c.Env)
		markedForRemoval += len(tweetsToDelete)

		numberRemoved, errs := deleteTweets(c, tweetsToDelete)
		if c.Env.Commit && c.Env.Verbose {
			fmt.Printf("\n")
		}
		c.Env.MaxAPICalls -= numberRemoved
		removed += numberRemoved
		errorCount += errs

		if errorCount < 20 && len(tweets) > 0 && c.Env.MaxAPICalls > 0 {
			max = tweets[len(tweets)-1].ID - 1
			if c.Env.Verbose {
				fmt.Printf("%v errs -- %v tweets -- %v calls left -- %v current id\n", errorCount, len(tweets), c.Env.MaxAPICalls, max)
			}
		} else {
			shouldContinue = false
		}
	}

	fmt.Printf("\nTotal Scanned Tweets: %v; Removed: %v of %v; Max Age: %v\n", count, removed, markedForRemoval, c.Env.MaxAge.Format(time.RFC3339))

	return nil
}
