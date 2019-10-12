package main

import (
	"fmt"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func isAgedOut(t *twitter.Tweet, env *PrunerEnv) bool {
	createdTime, _ := t.CreatedAtTime()

	return env.MaxAge.After(createdTime)
}

func isBoring(t *twitter.Tweet, env *PrunerEnv) bool {
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

func getTweetsToDelete(tweets []twitter.Tweet, env *PrunerEnv) []int64 {
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

func deleteTweet(te *twitter.Client, id int64) error {
	_, resp, err := te.Statuses.Destroy(id, &twitter.StatusDestroyParams{ID: id})
	if resp.StatusCode == 429 {
		wait, _ := time.ParseDuration(resp.Header.Get("x-rate-limit-reset") + "s")
		fmt.Printf("\nRate limit exceeded, waiting %v before trying again.\n", wait.String())
		<-time.After(wait)
		return deleteTweet(te, id)
	}
	if err != nil {
		return err
	}
	return nil
}

func deleteTweets(te *twitter.Client, tweetIds []int64, env *PrunerEnv) (int, int) {
	count := 0
	errorCount := 0
	if env.Commit {
		for _, id := range tweetIds {
			err := deleteTweet(te, id)
			if err != nil {
				if env.Verbose {
					fmt.Printf("\n")
				}
				fmt.Printf("Error removing status: %v\n", err)
				errorCount++
				continue
			}
			if env.Verbose {
				fmt.Printf(".")
			}
			count++
		}
	}
	return count, errorCount
}

// PruneTimeline does exactly what it says it does
func PruneTimeline(te *twitter.Client, user *twitter.User, env *PrunerEnv) error {
	count := 0
	markedForRemoval := 0
	removed := 0
	errorCount := 0
	boolTrue := true
	opts := &twitter.UserTimelineParams{Count: env.MaxTweetsPerRequest, TrimUser: &boolTrue, IncludeRetweets: &boolTrue}
	shouldContinue := true

	for shouldContinue {
		env.MaxAPICalls--
		tweets, _, err := te.Timelines.UserTimeline(opts)
		if err != nil {
			fmt.Printf("Error in timeline retrieval: %+v", err)
			errorCount++
		}
		count += len(tweets)

		tweetsToDelete := getTweetsToDelete(tweets, env)
		markedForRemoval += len(tweetsToDelete)

		numberRemoved, errs := deleteTweets(te, tweetsToDelete, env)
		if env.Commit && env.Verbose {
			fmt.Printf("\n")
		}
		env.MaxAPICalls -= numberRemoved
		removed += numberRemoved
		errorCount += errs

		if errorCount < 20 && len(tweets) > 0 && env.MaxAPICalls > 0 {
			opts.MaxID = tweets[len(tweets)-1].ID - 1
			if env.Verbose {
				fmt.Printf("%v errs -- %v tweets -- %v calls left -- %v current id\n", errorCount, len(tweets), env.MaxAPICalls, opts.MaxID)
			}
		} else {
			shouldContinue = false
		}
	}

	fmt.Printf("\nTotal Scanned Tweets: %v; Removed: %v of %v; Max Age: %v\n", count, removed, markedForRemoval, env.MaxAge.Format(time.RFC3339))

	return nil
}
