package main

import (
	"encoding/json"
	_ "fmt"
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

	var dictFile *os.File
	dictFile, err = os.Create("data/dict.json")
	if err != nil {
		log.Fatal(err)
	}
	defer dictFile.Close()

	ss, _ := update.StoreTweetsAndLatestId(tweetsFile, idFile, 12345)

	ms := make(map[string][]string)
	for _, tw := range ss {
		update.UpdateMarkovSpace(ms, tw)
	}
	var dict []byte
	dict, err = json.Marshal(ms)
	if err != nil {
		log.Fatal(err)
	}

	dictFile.Write(dict)
}
