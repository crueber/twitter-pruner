package main

import "github.com/dghubble/go-twitter/twitter"

const (
	maxLikesPerPage = 100
)

// PruneLikes does exactly what it says it does
func PruneLikes(te *twitter.Client, env *TwitterEnv) error {
	return nil
}

// The Ruby Code

// user = api_call :user, @options[:username]
// tweets_to_unlike = []
// tweets_to_delete = []

// puts "==> Checking likes..."
// total_likes = [user.favorites_count, MAX_API_TWEETS].min
// oldest_likes_page = (total_likes / MAX_LIKES_PER_PAGE).ceil

// oldest_likes_page.downto(1) do |page|
//   tweets = api_call :favorites, count: MAX_LIKES_PER_PAGE, page: page
//   tweets_to_unlike += tweets.reject(&method(:too_new?))
// end

// puts "==> Unliking #{tweets_to_unlike.size} tweets"
// tweets_not_found = []
// tweets_to_unlike.each_slice(MAX_TWEETS_PER_REQUEST) do |tweets|
//   begin
//     # api_call :unfavorite, tweets
//   rescue Twitter::Error::NotFound
//     tweets_not_found += tweets
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
