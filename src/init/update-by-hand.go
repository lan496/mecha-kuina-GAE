package main

import (
	"fmt"
	"github.com/lan496/mecha-kuina/src/update"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	idData, err := ioutil.ReadFile("data/latestId.csv")
	if err != nil {
		log.Fatal(err)
	}

	sinceId, _ := strconv.ParseInt(string(idData), 10, 64)
	fmt.Println("%s", string(idData))
	fmt.Println("sinceId", sinceId)

	var twf *os.File
	twf, err = os.OpenFile("data/tweets.csv", os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer twf.Close()

	var idf *os.File
	idf, err = os.Open("data/latestId.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer idf.Close()

	update.StoreTweetsAndLatestId(twf, idf, sinceId+1)
}
