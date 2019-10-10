package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

// Verify gets User baseline data.
func Verify(twit *twitter.Client, te *TwitterEnv) error {
	verifyParams := &twitter.AccountVerifyParams{SkipStatus: twitter.Bool(true), IncludeEmail: twitter.Bool(true)}
	user, _, err := twit.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return err
	}

	fmt.Printf("Verified :: %v :: %v :: %v\n", user.ScreenName, user.Name, user.Description)
	return nil
}
