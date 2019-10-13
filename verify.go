package main

import (
	"fmt"

	"github.com/crueber/twitter-pruner/pruner"
	"github.com/dghubble/go-twitter/twitter"
)

// Verify gets User baseline data.
func Verify(c *pruner.Client) (*twitter.User, error) {
	user, err := c.GetUserInfo()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Verified :: %v :: %v :: %v\n", user.ScreenName, user.Name, user.Description)
	fmt.Printf("%v statuses :: %v favorites\n", user.StatusesCount, user.FavouritesCount)
	return user, nil
}
