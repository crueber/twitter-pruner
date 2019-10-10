## Twitter Pruner

#### Still Under Construction

This tool is used to prune tweets on twitter. It was heavily inspired by MikeMcQuaid's project [Twitter Delete](https://github.com/MikeMcQuaid/TwitterDelete). None of his code was re-used, as this is written in Golang. Ruby is a good programming language, but if one wants to distribute binaries, it isn't ideal. Releases will be maintained for all major operating systems, so that people don't have to be familiar with a programming language to perform the task.

The features will eventually have (at least) parity with TwitterDelete.
* Delete old tweets and retweets with adjustable age preference.
* Unlike old tweets with adjustable age preference.
* Keep tweets with a particular like or retweet amount.

Uncertain on this one
* Delete tweets no longer exposed by Twitter API from a downloaded Twitter archive file


### Pre-reqs

You will need to get an app set up on Twitter in order to use this. Visit https://apps.twitter.com and take note of the following items, as you'll need them to run the pruner:

twitter_consumer_key
twitter_consumer_secret
twitter_access_token
twitter_access_token_secret

### Usage

Download an appropriate binary, and run it by typing `./twitter-pruner -h` in to the command console.

The basic command is this: `./twitter-pruner --key="<twitter_consumer_key>" --secret="<twitter_consumer_secret>" --token="<twitter_access_token>" --tsecret="<twitter_access_token_secret>"`

## License

Copyright 2019 by Christopher WJ Rueber

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.