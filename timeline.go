package main

import (
	"fmt"
	"os"
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
		fmt.Printf("Ignoring tweet (%v fav/%v rt): %v\n", t.FavoriteCount, t.RetweetCount, t.Text)
		return false
	}
	return true
}

func getTweetsToDelete(tweets []twitter.Tweet, env *PrunerEnv) []twitter.Tweet {
	var tweetsToDelete []twitter.Tweet
	for _, tweet := range tweets {
		if isAgedOut(&tweet, env) && isBoring(&tweet, env) {
			tweetsToDelete = append(tweetsToDelete, tweet)
		}
	}
	return tweetsToDelete
}

func deleteTweets(te *twitter.Client, tweets []twitter.Tweet, env *PrunerEnv) (int, int) {
	count := 0
	errorCount := 0
	if env.Commit {
		for _, tweet := range tweets {
			_, _, err := te.Statuses.Destroy(tweet.ID, &twitter.StatusDestroyParams{ID: tweet.ID})
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
	totalCount := 0
	totalMarkedForRemoval := 0
	totalRemoved := 0
	errorCount := 0
	inclRTs := true
	opts := &twitter.UserTimelineParams{Count: env.MaxTweetsPerRequest, IncludeRetweets: &inclRTs}
	shouldContinue := true
	if env.Commit {
		os.Exit(1)
	}

	for shouldContinue {
		tweets, _, err := te.Timelines.UserTimeline(opts)
		if err != nil {
			fmt.Printf("Error in timeline retrieval: %+v", err)
			errorCount++
		}

		tweetsToDelete := getTweetsToDelete(tweets, env)
		removed, errs := deleteTweets(te, tweetsToDelete, env)

		totalCount += len(tweets)
		totalMarkedForRemoval += len(tweetsToDelete)
		totalRemoved += removed
		errorCount += errs

		if errorCount < 20 && len(tweets) == env.MaxTweetsPerRequest && totalCount < env.MaxAPITweets {
			opts.MaxID = tweets[19].ID
			// fmt.Printf(".")
		} else {
			shouldContinue = false
		}
	}

	fmt.Printf("\nTotal Count: %v; Removed: %v of %v; Max Age: %v\n", totalCount, totalRemoved, totalMarkedForRemoval, env.MaxAge.Format(time.RFC3339))

	return nil
}
