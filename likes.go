package main

import (
	"fmt"
	"time"

	"github.com/crueber/twitter-pruner/pruner"
	"github.com/dghubble/go-twitter/twitter"
)

func whichTweetsToUnfavorite(tweets []twitter.Tweet, env *pruner.Env) []int64 {
	var unfav []int64

	for _, tweet := range tweets {
		if isAgedOut(&tweet, env) {
			unfav = append(unfav, tweet.ID)
			if env.Verbose {
				fmt.Printf("Unliking: %v --- %v\n", tweet.CreatedAt, tweet.Text)
			}
		}
	}

	return unfav
}

func processUnfavorite(c *pruner.Client, tweetIds []int64) (int, int) {
	count := 0
	errCount := 0
	if c.Env.Commit {
		for _, id := range tweetIds {
			err := c.DestroyLike(id)
			if err != nil {
				if c.Env.Verbose {
					fmt.Printf("\n")
				}
				fmt.Printf("%v\n", err)
				errCount++
				continue
			}
			if c.Env.Verbose {
				fmt.Printf(".")
			}
			count++
		}
	}
	return count, errCount
}

// PruneLikes does exactly what it says it does
func PruneLikes(c *pruner.Client, user *twitter.User) error {
	var max int64
	count := 0
	markedForRemoval := 0
	removed := 0
	errorCount := 0
	shouldContinue := true

	for shouldContinue {
		c.Env.MaxAPICalls--
		favs, err := c.GetLikes(max)
		if err != nil {
			fmt.Printf("Error retrieving favorites: %+v", err)
			errorCount++
			continue
		}

		unfav := whichTweetsToUnfavorite(favs, c.Env)
		unfaved, errs := processUnfavorite(c, unfav)
		if c.Env.Commit && c.Env.Verbose && unfaved > 0 {
			fmt.Printf("\n")
		}

		c.Env.MaxAPICalls -= unfaved
		count += len(favs)
		markedForRemoval += len(unfav)
		removed += unfaved
		errorCount += errs

		if errorCount < 20 && len(favs) > 0 && c.Env.MaxAPICalls > 0 {
			max = favs[len(favs)-1].ID - 1
			if c.Env.Verbose {
				fmt.Printf("%v errs -- %v likes -- %v calls left -- oldest in batch: %v\n", errorCount, len(favs), c.Env.MaxAPICalls, favs[len(favs)-1].CreatedAt)
			} else {
				fmt.Printf(".")
			}
		} else {
			if !c.Env.Verbose {
				fmt.Printf(".\n")
			}
			shouldContinue = false
		}
	}

	fmt.Printf("\nTotal Scanned Tweets: %v; Unliked: %v of %v; Max Age: %v\n", count, removed, markedForRemoval, c.Env.MaxAge.Format(time.RFC3339))

	return nil
}
