package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

const (
	maxAPITweets        = 3200
	maxTweetsPerPage    = 200
	maxTweetsPerRequest = 100
)

func getTweetsToDelete(te *twitter.Client, env *TwitterEnv) ([]twitter.Tweet, error) {
	tweetsToDelete := []twitter.Tweet{}

	return tweetsToDelete, nil
}

// PruneTimeline does exactly what it says it does
func PruneTimeline(te *twitter.Client, user *twitter.User, env *TwitterEnv) error {
	totalCount := 0

	// &twitter.UserTimelineParams{page: 1, count: 20}

	// Home Timeline
	// tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
	// 	Count: 20,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Printf("Total Count: %v\n", totalCount)

	return nil
}

// The Ruby Code

// puts "==> Checking timeline..."
// total_tweets = [user.statuses_count, MAX_API_TWEETS].min
// oldest_tweets_page = (total_tweets / MAX_TWEETS_PER_PAGE).ceil

// oldest_tweets_page.downto(1) do |page|
//   tweets = api_call :user_timeline, count: MAX_TWEETS_PER_PAGE, page: page
//   tweets_to_delete += tweets.reject(&method(:too_new_or_popular?))
// end

// puts "==> Deleting #{tweets_to_delete.size} tweets"
// tweets_to_delete.each_slice(MAX_TWEETS_PER_REQUEST) do |tweets|
//   begin
//     # api_call :destroy_status, tweets
//   rescue Twitter::Error::NotFound
//     tweets_not_found += tweets
//   end
// end

// tweets_not_found.each do |tweet|
//   begin
//     # api_call :destroy_status, tweet
//   rescue Twitter::Error::NotFound
//     nil
//   end
// end

// @oldest_tweet_time_to_keep = Time.now - @options[:days] * 24 * 60 * 60
// @newest_tweet_time_to_keep = Time.now - @options[:olds] * 24 * 60 * 60

// def too_new?(tweet)
//   tweet.created_at > @oldest_tweet_time_to_keep || tweet.created_at < @newest_tweet_time_to_keep
// end

// def too_new_or_popular?(tweet)
//   return true if too_new? tweet

//   return false if tweet.retweeted?
//   return false if tweet.text.start_with? "RT @"

//   if tweet.retweet_count >= @options[:rts]
//     puts "Ignoring tweet: too RTd: #{tweet.text}"
//     return true
//   end

//   if tweet.favorite_count >= @options[:favs]
//     puts "Ignoring tweet: too liked: #{tweet.text}"
//     return true
//   end

//   false
// end
