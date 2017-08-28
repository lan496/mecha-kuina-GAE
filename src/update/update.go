package update

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lan496/mecha-kuina/src/secret"
	"github.com/lan496/mecha-kuina/src/twitter"
	"os"
	"strconv"
	"strings"
)

func StoreTweetsAndLatestId(twf, idf *os.File) {
	anaconda.SetConsumerKey(secret.ConsumerKey)
	anaconda.SetConsumerSecret(secret.ConsumerSecret)
	api := anaconda.NewTwitterApi(secret.AccessToken, secret.AccessTokenSecret)

	tweets, latestId, _ := twitter.FetchTweets(api, secret.Username, 2000, 12345)

	fmt.Println("fetched tweets:", len(tweets))
	fmt.Println("latestId:", latestId)

	idf.WriteString(strconv.FormatInt(latestId, 10))

	for _, tw := range tweets {
		surfaces := twitter.Tokenize(twitter.TrimURL(tw))
		storedStr := strings.Join(surfaces, " ") + "\n"
		fmt.Println(storedStr)
		twf.WriteString(storedStr)
	}
}
