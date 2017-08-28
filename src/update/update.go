package update

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"os"
	"strconv"
)

func StoreTweetsAndLatestId(twf, idf *os.File, sinceId int64) (ss [][]string, latestId int64) {
	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	var tweets []string
	tweets, latestId, _ = twitter.FetchTweets(api, secret.Username, 2000, sinceId)

	fmt.Println("fetched tweets:", len(tweets))
	fmt.Println("latestId:", latestId)

	tmpid, _ := json.Marshal(strconv.FormatInt(latestId, 10))
	idf.Write(tmpid)

	for _, tw := range tweets {
		surfaces := twitter.Tokenize(twitter.TrimURL(tw))
		fmt.Println(surfaces)
		ss = append(ss, surfaces)
	}
	tmpss, _ := json.Marshal(ss)
	twf.Write(tmpss)

	return
}

func UpdateMarkovSpace(ms map[string][]string, surfaces []string) {
	for i := 0; i < len(surfaces)-3; i++ {
		key1 := surfaces[i] + "^" + surfaces[i+1]
		key2 := surfaces[i+1] + "^" + surfaces[i+2]
		ms[key1] = append(ms[key1], key2)
	}
}
