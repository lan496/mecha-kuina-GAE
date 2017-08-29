package main

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
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

func ReadLatestId(r *http.Request, filename string) (latestId int64) {
	ctx := appengine.NewContext(r)
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get client", err)
	}

	reader, err := client.Bucket(bucket).Object(filename).NewReader(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get reader", err)
	}
	defer reader.Close()

	fmt.Fscanf(reader, "%d", &latestId)
	return
}

func Update(w http.ResponseWriter, r *http.Request) {
	latestId := ReadLatestId(r, "latestId.txt")
	fmt.Fprintln(w, latestId)
}

func init() {
	FetchTweetsGAE := FetchTweetsHandler(secret.Username, 100, 12345)

	http.HandleFunc(secret.TweetQuery, Tweet)
	http.HandleFunc(secret.FetchQuery, FetchTweetsGAE)
	http.HandleFunc(secret.UpdateQuery, Update)
}
