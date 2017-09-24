package Rudebot

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

	Nouns, NumNouns, err = readInput(NounPath)
	if err != nil {
		log.Printf("failed to find noun file %s: %s\n", NounPath, err)
		return
	}

	Adjectives, NumAdjectives, err = readInput(AdjectivePath)
	if err != nil {
		log.Printf("failed to find adjective file %s: %s\n", AdjectivePath, err)
		return
	}

	Adverbs, NumAdverbs, err = readInput(AdverbPath)
	if err != nil {
		log.Printf("failed to find adverb file %s: %s\n", AdverbPath, err)
		return
	}

	Verbs, NumVerbs, err = readInput(VerbPath)
	if err != nil {
		log.Printf("failed to find verb file %s: %s\n", VerbPath, err)
		return
	}
}

//turns file into string array, array length, and possible error
func readInput(path string) ([]string, int, error) {
	data, err := ioutil.ReadFile(path)
	str, num := sanitizeInput(data)
	return str, num, err
}

//turns byte array into string array and string array length
func sanitizeInput(data []byte) ([]string, int) {
	dataString := string(data)
	dataLines := strings.Split(dataString, "\n")

	cleanInput(dataLines)

	return dataLines, len(dataLines)
}

//Clean all white space for each string in array
func cleanInput(input []string) {
	for i, word := range input {
		input[i] = stripWhiteSpace(word)
	}
}

//Clean all white space from string and return it
func stripWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func startsWithVowel(str string) bool {
	chr := str[0]
	switch chr {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}

//sets LastAdjective to adj and LastNoun to noun
func saveInsult(adj string, noun string) {
	LastAdjective = adj
	LastNoun = noun
}

//Creates and insult directed at target, using adj and noun Stores adj and noun for rating
func createInsult(target string, adj string, noun string) string {
	_, ok := NounRatings[noun]
	if !ok {
		NounRatings[noun] = 0
	}

	_, ok = AdjectiveRatings[adj]
	if !ok {
		AdjectiveRatings[adj] = 0
	}

	saveInsult(adj, noun)

	var insult string

	if startsWithVowel(adj) {
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
	return createInsult(target, Adjectives[rand.Intn(NumAdjectives)], Nouns[rand.Intn(NumNouns)])
}

//Rates last used adjective and noun, changing rating by adding given value
func Rate(value int) {
	NounRatings[LastNoun] += value
	AdjectiveRatings[LastAdjective] += value
}

func splitPositive(words PairList) PairList {
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

func splitNegative(words PairList) PairList {
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

	adjectives, nouns := getRatingList()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	nouns = splitPositive(nouns)
	adjectives = splitPositive(adjectives)

	if len(nouns) == 0 || len(adjectives) == 0 {
		return "No good insults"
	}

	adj := adjectives[rand.Intn(len(adjectives))].Key
	noun := nouns[rand.Intn(len(nouns))].Key

	return createInsult(target, adj, noun)
}

func BadInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingList()

	sort.Sort(nouns)
	sort.Sort(adjectives)

	nouns = splitNegative(nouns)
	adjectives = splitNegative(adjectives)

	if len(nouns) == 0 || len(adjectives) == 0 {
		return "No bad insults"
	}

	adj := adjectives[rand.Intn(len(adjectives))].Key
	noun := nouns[rand.Intn(len(nouns))].Key

	return createInsult(target, adj, noun)
}

func BestInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingList()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	adj := adjectives[0].Key
	noun := nouns[0].Key

	return createInsult(target, adj, noun)
}

func WorstInsult(target string) string {
	if len(NounRatings) == 0 || len(AdjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingList()

	sort.Sort(nouns)
	sort.Sort(adjectives)

	adj := adjectives[0].Key
	noun := nouns[0].Key

	return createInsult(target, adj, noun)
}

func getRatingList() (PairList, PairList) {
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
	return createInsult(target, LastAdjective, LastNoun)
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
