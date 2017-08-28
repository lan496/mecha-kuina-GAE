package main

import (
	"github.com/lan496/mecha-kuina/src/update"
	"log"
	"os"
)

func main() {
	idFile, err := os.Create("data/latestId.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer idFile.Close()

	var tweetsFile *os.File
	tweetsFile, err = os.Create("data/tweets.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer tweetsFile.Close()

	update.StoreTweetsAndLatestId(tweetsFile, idFile)
}
