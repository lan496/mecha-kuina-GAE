package update

import (
	"bufio"
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"net/http"
	"strings"
)

func LatestTweetsAndId(sinceId int64) (ss []string, latestId int64, err error) {
	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	var tweets []string
	tweets, latestId, err = twitter.FetchTweets(api, secret.Username, 2000, sinceId)

	for _, tw := range tweets {
		surfaces := twitter.Tokenize(twitter.TrimURL(tw))
		storedStr := strings.Join(surfaces, " ") + "\n"
		fmt.Print(storedStr)
		ss = append(ss, storedStr)
	}
	return
}

func ReadLatestId(r *http.Request, filename string) (latestId int64, err error) {
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

func ReadTweets(r *http.Request, filename string) (tweets []string, err error) {
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

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		tweets = append(tweets, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Errorf(ctx, "failed to scan file", err)
	}
	return
}
