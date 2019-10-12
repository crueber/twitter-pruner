package main

import (
	"fmt"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func whichTweetsToUnfavorite(tweets []twitter.Tweet, env *PrunerEnv) []int64 {
	var unfav []int64

	for _, tweet := range tweets {
		if isAgedOut(&tweet, env) {
			unfav = append(unfav, tweet.ID)
		}
	}

	return unfav
}

func processUnfavorite(te *twitter.Client, env *PrunerEnv, tweetIds []int64) (int, int) {
	count := 0
	errCount := 0
	if env.Commit {
		for _, id := range tweetIds {
			_, _, err := te.Favorites.Destroy(&twitter.FavoriteDestroyParams{ID: id})
			if err != nil {
				fmt.Printf("Error removing favorite: %v", err)
				errCount++
			}
		}
	}
	return count, errCount
}

// PruneLikes does exactly what it says it does
func PruneLikes(te *twitter.Client, user *twitter.User, env *PrunerEnv) error {
	count := 0
	markedForRemoval := 0
	removed := 0
	errorCount := 0
	opts := &twitter.FavoriteListParams{Count: env.MaxTweetsPerRequest}
	shouldContinue := true

	for shouldContinue {
		favs, _, err := te.Favorites.List(opts)
		if err != nil {
			fmt.Printf("Error retrieving favorites: %+v", err)
			errorCount++
		}

		unfav := whichTweetsToUnfavorite(favs, env)
		unfaved, errs := processUnfavorite(te, env, unfav)

		env.MaxAPITweets -= len(favs)
		count += len(favs)
		markedForRemoval += len(unfav)
		removed += unfaved
		errorCount += errs

		if errorCount < 20 && len(favs) > 1 && env.MaxAPITweets > 0 {
			opts.MaxID = favs[len(favs)-1].ID
		} else {
			shouldContinue = false
		}
	}

	fmt.Printf("\nTotal Count: %v; Removed: %v of %v; Max Age: %v\n", count, removed, markedForRemoval, env.MaxAge.Format(time.RFC3339))

	return nil
}
