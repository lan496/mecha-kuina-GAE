package main

import (
	"fmt"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"os"
)

func main() {
	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	tweets, latestId, _ := twitter.FetchTweets(api, secret.Username, 2000, 12345)

	fmt.Println("fetched tweets:", len(tweets))
}
