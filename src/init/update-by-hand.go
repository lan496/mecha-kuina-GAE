package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/update"
	"log"
	"os"
)

func readLatestId(filename string) (latestId int64, err error) {
	idReader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer idReader.Close()

	fmt.Fscanf(idReader, "%d", &latestId)
	fmt.Println(latestId)

	return
}

func main() {
	sinceId, _ := readLatestId("data/latestId.txt")

	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	tweets, latestId, err := update.LatestTweetsAndId(api, sinceId)

	for _, tw := range tweets {
		fmt.Println(tw)
	}
	fmt.Println(latestId)

	tweetsFile, err := os.OpenFile("data/tweets.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer tweetsFile.Close()

	latestIdFile, err := os.OpenFile("data/latestId.txt", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer latestIdFile.Close()

	for _, tw := range tweets {
		fmt.Fprintf(tweetsFile, "%s\n", tw)
	}

	fmt.Fprintf(latestIdFile, "%d\n", latestId)
}
