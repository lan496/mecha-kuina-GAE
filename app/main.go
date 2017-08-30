package main

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/markov"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"github.com/lan496/mecha-kuina/src/update"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"net/url"
	"strings"
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

func Update(w http.ResponseWriter, r *http.Request) {
	prevLatestId, _ := update.ReadLatestId(r, "latestId.txt")
	fmt.Fprintln(w, prevLatestId)
	prevTweets, _ := update.ReadTweets(r, "tweets.txt")
	fmt.Fprintf(w, "%d tweets are collected.\n", len(prevTweets))

	tweets, latestId, _ := update.LatestTweetsAndId(prevLatestId)
	fmt.Fprintf(w, "%d tweets are collected anew.\n", len(tweets))
	fmt.Fprintln(w, latestId)

	ctx := appengine.NewContext(r)
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get client", err)
	}

	tweetsWriter := client.Bucket(bucket).Object("tweets.txt").NewWriter(ctx)
	defer tweetsWriter.Close()

	for _, tw := range tweets {
		fmt.Fprintf(tweetsWriter, "%s\n", tw)
	}
	for _, tw := range prevTweets {
		fmt.Fprintf(tweetsWriter, "%s\n", tw)
	}

	idWriter := client.Bucket(bucket).Object("latestId.txt").NewWriter(ctx)
	defer idWriter.Close()

	fmt.Fprintf(idWriter, "%d\n", latestId)
}

func Markov(w http.ResponseWriter, r *http.Request) {
	tweets, _ := update.ReadTweets(r, "tweets.txt")
	ms := make(markov.MarkovSpace2)

	for _, tw := range tweets {
		surfaces := strings.Split(tw, " ")
		markov.UpdateMarkovSpace2(ms, surfaces)
	}

	sentenceNum := 10
	for i := 0; i < sentenceNum; i++ {
		s := markov.CreateSentence2(ms)
		fmt.Fprintf(w, "%s\n", s)
	}
}

func init() {
	FetchTweetsGAE := FetchTweetsHandler(secret.Username, 100, 12345)

	http.HandleFunc(secret.TweetQuery, Tweet)
	http.HandleFunc(secret.FetchQuery, FetchTweetsGAE)
	http.HandleFunc(secret.UpdateQuery, Update)
	http.HandleFunc(secret.MarkovQuery, Markov)
}
