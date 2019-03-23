package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/update"
	"log"
	"os"
)

func main() {
	idFile, err := os.Create("data/latestId.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer idFile.Close()

	var tweetsFile *os.File
	tweetsFile, err = os.Create("data/tweets.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer tweetsFile.Close()

	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	tweets, latestId, _ := update.LatestTweetsAndId(api, 12345)

	fmt.Println(latestId)
	fmt.Fprintf(idFile, "%d\n", latestId)

	for _, tw := range tweets {
		fmt.Fprintf(tweetsFile, "%s\n", tw)
	}
}
