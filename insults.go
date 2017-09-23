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
	WordDirectory = "words"
	NounPath      = "words/nouns/all.txt"
	AdjectivePath = "words/adjectives/all.txt"
	AdverbPath    = "words/adverbs/all.txt"
	VerbPath      = "words/verbs/all.txt"
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

func InitRatings() {
	NounRatings = map[string]int{}
	AdjectiveRatings = map[string]int{}
	AdverbRatings = map[string]int{}
	VerbRatings = map[string]int{}
}

func InitInsults() {
	rand.Seed(time.Now().Unix())
	var err error

	Nouns, NumNouns, err = ReadInput(NounPath)
	if err != nil {
		log.Printf("failed to find noun file %s: %s\n", NounPath, err)
		return
	}

	Adjectives, NumAdjectives, err = ReadInput(AdjectivePath)
	if err != nil {
		log.Printf("failed to find adjective file %s: %s\n", AdjectivePath, err)
		return
	}

	Adverbs, NumAdverbs, err = ReadInput(AdverbPath)
	if err != nil {
		log.Printf("failed to find adverb file %s: %s\n", AdverbPath, err)
		return
	}

	Verbs, NumVerbs, err = ReadInput(VerbPath)
	if err != nil {
		log.Printf("failed to find verb file %s: %s\n", VerbPath, err)
		return
	}
}

//turns file into string array, array length, and possible error
func ReadInput(path string) ([]string, int, error) {
	data, err := ioutil.ReadFile(path)
	str, num := SanitizeInput(data)
	return str, num, err
}

//turns byte array into string array and string array length
func SanitizeInput(data []byte) ([]string, int) {
	dataString := string(data)
	dataLines := strings.Split(dataString, "\n")

	CleanInput(dataLines)

	return dataLines, len(dataLines)
}

//Clean all white space for each string in array
func CleanInput(input []string) {
	for i, word := range input {
		input[i] = StripWhiteSpace(word)
	}
}

//Clean all white space from string and return it
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

//sets LastAdjective to adj and LastNoun to noun
func SaveInsult(adj string, noun string) {
	LastAdjective = adj
	LastNoun = noun
}

//Creates and insult directed at target, using adj and noun Stores adj and noun for rating
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
	if len(Nouns) == 0 || len(Adjectives) == 0 {
		return "Not enough valid words"
	}
	return CreateInsult(target, Adjectives[RandomInt(NumAdjectives)], Nouns[RandomInt(NumNouns)])
}

func RandomInt(max int) int {
	return rand.Intn(max)
}

//Rates last used adjective and noun, changing rating by adding given value
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

	sort.Sort(nouns)
	sort.Sort(adjectives)

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

func LastInsult(target string) string {
	if LastNoun == "" || LastAdjective == "" {
		return "No previous insult"
	}
	return CreateInsult(target, LastAdjective, LastNoun)
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
