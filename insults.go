package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

const (
	wordDirectory = "words"
	nounPath      = "words/nouns/all.txt"
	adjectivePath = "words/adjectives/all.txt"
	adverbPath    = "words/adverbs/all.txt"
	verbPath      = "words/verbs/all.txt"
)

var (
	LastNoun      string
	LastAdjective string
	LastVerb      string
	LastAdverb    string

	NumAdjectives    int
	Adjectives       []string
	AdjectiveRatings map[string]int
	NumAdverbs       int
	Adverbs          []string
	AdverbRatings    map[string]int
	NumNouns         int
	Nouns            []string
	NounRatings      map[string]int
	NumVerbs         int
	Verbs            []string
	VerbRatings      map[string]int
)

func InitInsults() {

	rand.Seed(time.Now().Unix())

	NounRatings = map[string]int{}
	AdjectiveRatings = map[string]int{}
	AdverbRatings = map[string]int{}
	VerbRatings = map[string]int{}

	nounData, err := ioutil.ReadFile(nounPath)
	if err != nil {
		log.Printf("failed to find noun file %s: %s\n", nounPath, err)
		return
	}
	dataString := string(nounData)
	dataLines := strings.Split(dataString, "\n")
	NumNouns = len(dataLines)
	Nouns = dataLines

	CleanInput(Nouns)

	adjectiveData, err := ioutil.ReadFile(adjectivePath)
	if err != nil {
		log.Printf("failed to find noun file %s: %s\n", adjectivePath, err)
		return
	}

	dataString = string(adjectiveData)
	dataLines = strings.Split(dataString, "\n")
	NumAdjectives = len(dataLines)
	Adjectives = dataLines

	CleanInput(Adjectives)

	adverbData, err := ioutil.ReadFile(adverbPath)
	if err != nil {
		log.Printf("failed to find noun file %s: %s\n", adverbPath, err)
		return
	}

	dataString = string(adverbData)
	dataLines = strings.Split(dataString, "\n")
	NumAdverbs = len(dataLines)

	verbData, err := ioutil.ReadFile(verbPath)
	if err != nil {
		log.Printf("failed to find noun file %s: %s\n", verbPath, err)
		return
	}

	dataString = string(verbData)
	dataLines = strings.Split(dataString, "\n")
	NumVerbs = len(dataLines)

}

func CleanInput(input []string) {
	for i, word := range input {
		input[i] = StripWhiteSpace(word)
	}
}

func StripWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func StartsWithVowel(str string) bool {
	chr := str[0]
	switch chr {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}

func SaveInsult(adj string, noun string) {
	LastAdjective = adj
	LastNoun = noun
}

func CreateInsult(target string, adj string, noun string) string {
	_, ok := NounRatings[noun]
	if !ok {
		NounRatings[noun] = 0
	}

	_, ok = AdjectiveRatings[adj]
	if !ok {
		AdjectiveRatings[adj] = 0
	}

	SaveInsult(adj, noun)

	var insult string

	if StartsWithVowel(adj) {
		insult = fmt.Sprintf("%s is an %s %s (%d,%d)", target, adj, noun, AdjectiveRatings[adj], NounRatings[noun])
	} else {
		insult = fmt.Sprintf("%s is a %s %s (%d,%d)", target, adj, noun, AdjectiveRatings[adj], NounRatings[noun])
	}

	return insult
}

func RandomInsult(target string) string {
	return CreateInsult(target, Adjectives[RandomInt(NumAdjectives)], Nouns[RandomInt(NumNouns)])
}

func RandomInt(max int) int {
	return rand.Intn(max)
}

func Rate(value int) {
	NounRatings[LastNoun] += value
	AdjectiveRatings[LastAdjective] += value
}

func SplitPositive(words PairList) PairList {
	var lastZero int
	for k, v := range words {
		if v.Value >= 0 {
			lastZero = k + 1
		} else {
			break
		}
	}
	newWords := words[0:lastZero]

	return newWords
}

func SplitNegative(words PairList) PairList {
	var lastZero int
	for k, v := range words {
		if v.Value <= 0 {
			lastZero = k + 1
		} else {
			break
		}
	}
	newWords := words[0:lastZero]

	return newWords
}

func GoodInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := GetRatingLists()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	nouns = SplitPositive(nouns)
	adjectives = SplitPositive(adjectives)

	if len(nouns) == 0 || len(adjectives) == 0 {
		return "No good insults"
	}

	adj := adjectives[RandomInt(len(adjectives))].Key
	noun := nouns[RandomInt(len(nouns))].Key

	return CreateInsult(target, adj, noun)
}

func BadInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := GetRatingLists()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	nouns = SplitNegative(nouns)
	adjectives = SplitNegative(adjectives)

	if len(nouns) == 0 || len(adjectives) == 0 {
		return "No bad insults"
	}

	adj := adjectives[RandomInt(len(adjectives))].Key
	noun := nouns[RandomInt(len(nouns))].Key

	return CreateInsult(target, adj, noun)
}

func BestInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := GetRatingLists()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	adj := adjectives[0].Key
	noun := nouns[0].Key

	return CreateInsult(target, adj, noun)
}

func WorstInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := GetRatingLists()

	sort.Sort(nouns)
	sort.Sort(adjectives)

	adj := adjectives[0].Key
	noun := nouns[0].Key

	return CreateInsult(target, adj, noun)
}

func GetRatingLists() (PairList, PairList) {
	var nouns PairList
	for k, v := range NounRatings {
		nouns = append(nouns, Pair{k, v})
	}

	var adjectives PairList
	for k, v := range AdjectiveRatings {
		adjectives = append(adjectives, Pair{k, v})
	}
	sort.Sort(sort.Reverse(adjectives))

	return adjectives, nouns
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
