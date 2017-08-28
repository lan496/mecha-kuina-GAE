package twitter

import (
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"github.com/ikawaha/kagome.ipadic/tokenizer"
	"net/url"
	"regexp"
	"strconv"
)

func FetchTweets(a *anaconda.TwitterApi, username string, n int, sinceId int64) (tweets []string, latestId int64, err error) {
	tweets = make([]string, 0)
	latestId = 0

	count := 200
	gettimes := (n + count - 1) / count
	var maxId int64

	v := url.Values{}
	v.Set("screen_name", username)
	v.Set("exclude_replies", "true")
	v.Set("include_rts", "false")
	v.Set("since_id", strconv.FormatInt(sinceId, 10))
	v.Set("count", strconv.Itoa(count))

	for i := 0; i < gettimes; i++ {
		if maxId != 0 {
			v.Set("max_id", strconv.FormatInt(maxId, 10))
		}

		var timeline []anaconda.Tweet
		timeline, err = a.GetUserTimeline(v)

		if err != nil {
			return
		}
		if len(timeline) == 0 {
			err = errors.New("cannot fetch tweets any more.")
			return
		}

		for _, tw := range timeline {
			tweets = append(tweets, tw.Text)
		}

		if maxId != 0 {
			latestId = timeline[0].Id
		}
		maxId = timeline[len(timeline)-1].Id - 1
	}
	return
}

func TrimURL(s string) (r string) {
	rep := regexp.MustCompile("https?://[0-9A-Za-z/:%#$&?()~.=+-]+")
	r = rep.ReplaceAllString(s, "")
	return
}

func Tokenize(tweet string) (surfaces []string) {
	t := tokenizer.New()
	tokens := t.Tokenize(tweet)
	surfaces = make([]string, 0)

	for _, token := range tokens {
		surfaces = append(surfaces, token.Surface)
	}
	return
}