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
		fmt.Printf("Ignoring %vtweet (%v fav/%v rt): %v\n", rt, t.FavoriteCount, t.RetweetCount, t.Text)
		return false
	}
	return true
}

func getTweetsToDelete(tweets []twitter.Tweet, env *PrunerEnv) []int64 {
	var tweetsToDelete []int64
	for _, tweet := range tweets {
		if isAgedOut(&tweet, env) && isBoring(&tweet, env) {
			tweetsToDelete = append(tweetsToDelete, tweet.ID)
		}
	}
	return tweetsToDelete
}

func deleteTweets(te *twitter.Client, tweetIds []int64, env *PrunerEnv) (int, int) {
	count := 0
	errorCount := 0
	if env.Commit {
		for _, id := range tweetIds {
			_, _, err := te.Statuses.Destroy(id, &twitter.StatusDestroyParams{ID: id})
			if err != nil {
				fmt.Printf("Error removing status: %v", err)
				errorCount++
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
		tweets, _, err := te.Timelines.UserTimeline(opts)
		if err != nil {
			fmt.Printf("Error in timeline retrieval: %+v", err)
			errorCount++
		}

		tweetsToDelete := getTweetsToDelete(tweets, env)
		removed, errs := deleteTweets(te, tweetsToDelete, env)

		env.MaxAPITweets -= len(tweets)
		count += len(tweets)
		markedForRemoval += len(tweetsToDelete)
		removed += removed
		errorCount += errs

		if errorCount < 20 && len(tweets) > 1 && env.MaxAPITweets > 0 {
			opts.MaxID = tweets[len(tweets)-1].ID
			fmt.Printf("%v -- %v -- %v -- %v", errorCount, len(tweets), env.MaxAPITweets, opts.MaxID)
		} else {
			shouldContinue = false
		}
	}

	fmt.Printf("\nTotal Count: %v; Removed: %v of %v; Max Age: %v\n", count, removed, markedForRemoval, env.MaxAge.Format(time.RFC3339))

	return nil
}
