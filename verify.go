package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

// Verify gets User baseline data.
func Verify(twit *twitter.Client, te *PrunerEnv) (*twitter.User, error) {
	verifyParams := &twitter.AccountVerifyParams{SkipStatus: twitter.Bool(true), IncludeEmail: twitter.Bool(true)}
	user, _, err := twit.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Verified :: %v :: %v :: %v\n", user.ScreenName, user.Name, user.Description)
	fmt.Printf("Total number of statuses: %v\n", user.StatusesCount)
	return user, nil
}
