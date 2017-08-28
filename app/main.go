package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"net/url"
)

func FetchTweetsHandler(username string, n int, sinceId int64) func(http.ResponseWriter, *http.Request) {
	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		api.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

		tweets, _, _ := twitter.FetchTweets(api, username, n, sinceId)
		for _, tw := range tweets {
			surfaces := twitter.Tokenize(twitter.TrimURL(tw))
			fmt.Fprint(w, surfaces)
		}
	}
}

func Tweet(w http.ResponseWriter, r *http.Request) {
	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)
	ctx := appengine.NewContext(r)
	api.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

	v := url.Values{}
	_, err := api.PostTweet("ファンタジーは好き", v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, "kuinakuina")
}

func init() {
	FetchTweetsGAE := FetchTweetsHandler(secret.Username, 100, 12345)

	http.HandleFunc(secret.TweetQuery, Tweet)
	http.HandleFunc(secret.FetchQuery, FetchTweetsGAE)
}
