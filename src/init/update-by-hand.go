package main

import (
	"fmt"
	"github.com/lan496/mecha-kuina/src/update"
	"log"
	"os"
)

func main() {
	idReader, err := os.Open("data/latestId.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer idReader.Close()

	var sinceId int64
	fmt.Fscanf(idReader, "%d", &sinceId)
	fmt.Println(sinceId)

	_, latestId, err := update.LatestTweetsAndId(sinceId)

	/*
		for _, tw := range tweets {
			fmt.Println(tw)
		}
	*/

	fmt.Println(latestId)

}
