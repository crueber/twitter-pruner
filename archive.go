package main

import "github.com/dghubble/go-twitter/twitter"

// PruneArchive prunes the twitter archive
func PruneArchive(te *twitter.Client, user *twitter.User, env *PrunerEnv) error {
	return nil
}

// The Ruby Code

// if @options[:archive_given]
//   puts "==> Checking archive JS..."
//   archive_tweet_ids = []

//   # tweet.js is not valid JSON...
//   file_contents = File.read(@options[:archive])
//   file_contents.sub! 'window.YTD.tweet.part0 = ', ''

//   JSON.parse(file_contents).each do |tweet|
//     archive_tweet_ids << tweet["id_str"]
//   end

//   archive_tweet_ids.each_slice(MAX_TWEETS_PER_REQUEST) do |tweet_ids|
//     tweets = api_call :statuses, tweet_ids
//     tweets_to_delete += tweets.reject(&method(:too_new_or_popular?))
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
