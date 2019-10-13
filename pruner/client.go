package pruner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// Client is the wrapper for a twitter client that is used by all the calls in this package
type Client struct {
	T   *twitter.Client
	Env *Env
}

func isRateLimited(resp *http.Response) bool {
	if resp.StatusCode == 429 {
		wait, _ := time.ParseDuration(resp.Header.Get("x-rate-limit-reset") + "s")
		fmt.Printf("Rate limit exceeded, waiting %v before trying again.", wait.String())
		<-time.After(wait)
		return true
	}
	return false
}

// DestroyLike removes a twitter like based on its ID
func (c *Client) DestroyLike(id int64) error {
	_, resp, err := c.T.Favorites.Destroy(&twitter.FavoriteDestroyParams{ID: id})
	if err != nil {
		if isRateLimited(resp) {
			return c.DestroyLike(id)
		}
		return err
	}
	return nil
}

// DestroyTweet removes a tweet based on its ID
func (c *Client) DestroyTweet(id int64) error {
	_, resp, err := c.T.Statuses.Destroy(id, &twitter.StatusDestroyParams{ID: id})
	if err != nil {
		if isRateLimited(resp) {
			return c.DestroyTweet(id)
		}
		return err
	}
	return nil
}

// GetTimeline gets timeline based on an empty max or a max identified
func (c *Client) GetTimeline(max int64) ([]twitter.Tweet, error) {
	opts := &twitter.UserTimelineParams{Count: c.Env.MaxTweetsPerRequest, TrimUser: twitter.Bool(true), IncludeRetweets: twitter.Bool(true)}

	if max > 0 {
		opts.MaxID = max
	}

	tweets, resp, err := c.T.Timelines.UserTimeline(opts)
	if err != nil {
		if isRateLimited(resp) {
			return c.GetTimeline(max)
		}
		return nil, err
	}

	return tweets, nil
}

// GetLikes gets likes based on an empty max or a max identified
func (c *Client) GetLikes(max int64) ([]twitter.Tweet, error) {
	opts := &twitter.FavoriteListParams{Count: c.Env.MaxTweetsPerRequest, IncludeEntities: twitter.Bool(false)}

	if max > 0 {
		opts.MaxID = max
	}

	favs, resp, err := c.T.Favorites.List(opts)
	if err != nil {
		if isRateLimited(resp) {
			return c.GetLikes(max)
		}
		return nil, err
	}

	return favs, nil
}

// GetUserInfo returns details needed to begin pulling data
func (c *Client) GetUserInfo() (*twitter.User, error) {
	verifyParams := &twitter.AccountVerifyParams{SkipStatus: twitter.Bool(true), IncludeEmail: twitter.Bool(false)}
	user, _, err := c.T.Accounts.VerifyCredentials(verifyParams)
	return user, err
}
