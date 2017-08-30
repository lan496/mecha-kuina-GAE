package markov

import (
	"math/rand"
	"time"
)

type MarkovSpace map[[3]string][][3]string
type MarkovSpace2 map[[2]string][][2]string

func UpdateMarkovSpace(ms MarkovSpace, surfaces []string) {
	if len(surfaces) < 2 {
		return
	}
	w1, w2, w3 := "", "", surfaces[0]
	for i := 0; i < len(surfaces)-1; i++ {
		key := [3]string{w1, w2, w3}
		value := [3]string{w2, w3, surfaces[i+1]}
		ms[key] = append(ms[key], value)
		w1, w2, w3 = w2, w3, surfaces[i+1]
	}
}

func UpdateMarkovSpace2(ms MarkovSpace2, surfaces []string) {
	if len(surfaces) < 2 {
		return
	}
	w1, w2 := "", surfaces[0]
	for i := 0; i < len(surfaces)-1; i++ {
		key := [2]string{w1, w2}
		value := [2]string{w2, surfaces[i+1]}
		ms[key] = append(ms[key], value)
		w1, w2 = w2, surfaces[i+1]
	}
}

func CreateSentence2(ms MarkovSpace2) (s string) {
	rand.Seed(time.Now().UnixNano())
	key := [2]string{"", "BOS"}
	for {
		if _, ok := ms[key]; ok == false {
			break
		}
		randomIdx := rand.Intn(len(ms[key]))
		value := ms[key][randomIdx]

		if value[1] == "EOS" {
			break
		}
		s += value[1]
		key = value
	}
	return
}
