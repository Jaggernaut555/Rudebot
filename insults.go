package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

type pair struct {
	Key   string
	Value int
}

type pairList []pair

const (
	wordDirectory       = "words"
	nounPath            = "words/nouns/all.txt"
	adjectivePath       = "words/adjectives/all.txt"
	adverbPath          = "words/adverbs/all.txt"
	verbPath            = "words/verbs/all.txt"
	ratingPath          = "ratings"
	nounRatingPath      = "ratings/nouns"
	adjectiveRatingPath = "ratings/adjectives"
)

var (
	lastNoun      map[string]string
	lastAdjective map[string]string

	NumAdjectives    int
	Adjectives       []string
	adjectiveRatings map[string]int
	NumAdverbs       int
	Adverbs          []string
	NumNouns         int
	Nouns            []string
	nounRatings      map[string]int
	NumVerbs         int
	Verbs            []string
)

func InitRatings() {
	nounRatings = make(map[string]int)
	adjectiveRatings = make(map[string]int)

	// If rating path does not exist create it
	if _, err := os.Stat(ratingPath); os.IsNotExist(err) {
		err = os.Mkdir(ratingPath, 0755)
		if err != nil {
			log.Printf("Error creating directory: %s\n", err)
		}
		fmt.Printf("Creating ratings directory %s\n", ratingPath)
	}

	loadRatings(&nounRatings, nounRatingPath)
	loadRatings(&adjectiveRatings, adjectiveRatingPath)
}

// Load the given rating map from the path
func loadRatings(list *map[string]int, path string) {
	ratingsFile, err := os.Open(path)
	defer ratingsFile.Close()
	if os.IsNotExist(err) {
		fmt.Printf("No file %s\n", path)
	} else {
		decoder := gob.NewDecoder(ratingsFile)
		err = decoder.Decode(list)
		if err != nil {
			log.Printf("Error loading file %s: %s\n", path, err)
		}
		fmt.Printf("loaded %d ratings from %s\n", len(*list), path)
	}
}

// save all ratings maps
func saveRatings() {
	saveMap(nounRatings, nounRatingPath)
	saveMap(adjectiveRatings, adjectiveRatingPath)
}

// save the ratings map to given path
func saveMap(data map[string]int, path string) {
	ratingFile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	defer ratingFile.Close()
	if err != nil {
		log.Printf("failed to open or create file %s: %s\n", path, err)
	}

	encoder := gob.NewEncoder(ratingFile)

	if err := encoder.Encode(data); err != nil {
		log.Printf("Failed to save ratings to %s: %s", path, err)
	}
}

//Rates last used adjective and noun, changing rating by adding given value
func Rate(channel string, value int) {
	nounRatings[lastNoun[channel]] += value
	adjectiveRatings[lastAdjective[channel]] += value

	saveRatings()
}

func InitInsults() {
	rand.Seed(time.Now().Unix())

	loadWords(nounPath, &Nouns, &NumNouns)
	loadWords(adjectivePath, &Adjectives, &NumAdjectives)
	loadWords(adverbPath, &Adverbs, &NumAdverbs)
	loadWords(verbPath, &Verbs, &NumVerbs)

	lastAdjective = map[string]string{}
	lastNoun = map[string]string{}
}

func loadWords(path string, words *[]string, num *int) {
	list, count, err := readWords(path)
	if err != nil {
		log.Printf("Failed to find file %s: %s\n", path, err)
	}
	*num = count
	*words = list
}

//turns file into string array, array length, and possible error
func readWords(path string) ([]string, int, error) {
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

//sets lastAdjective to adj and lastNoun to noun
func saveInsult(channel string, adj string, noun string) {
	lastAdjective[channel] = adj
	lastNoun[channel] = noun
}

//Creates and insult directed at target, using adj and noun Stores adj and noun for rating
func createInsult(channel string, target string, adj string, noun string) (insult string) {
	saveInsult(channel, adj, noun)
	Rate(channel, 0)

	if startsWithVowel(adj) {
		insult = fmt.Sprintf("%s is an %s %s (%d,%d)", target, adj, noun, adjectiveRatings[adj], nounRatings[noun])
	} else {
		insult = fmt.Sprintf("%s is a %s %s (%d,%d)", target, adj, noun, adjectiveRatings[adj], nounRatings[noun])
	}

	return insult
}

func RandomInsult(channel string, target string) string {
	if len(Nouns) == 0 || len(Adjectives) == 0 {
		return "Not enough valid words"
	}
	return createInsult(channel, target, Adjectives[rand.Intn(NumAdjectives)], Nouns[rand.Intn(NumNouns)])
}

func LastInsult(channel string, target string) string {
	if lastNoun[channel] == "" || lastAdjective[channel] == "" {
		return "No previous insult"
	}
	return createInsult(channel, target, lastAdjective[channel], lastNoun[channel])
}

func splitPositive(words pairList) pairList {
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

func splitNegative(words pairList) pairList {
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

func GoodInsult(channel string, target string) string {
	if len(nounRatings) == 0 || len(adjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingsLists()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	nouns = splitPositive(nouns)
	adjectives = splitPositive(adjectives)

	if len(nouns) == 0 || len(adjectives) == 0 {
		return "No good insults"
	}

	adj := adjectives[rand.Intn(len(adjectives))].Key
	noun := nouns[rand.Intn(len(nouns))].Key

	return createInsult(channel, target, adj, noun)
}

func BadInsult(channel string, target string) string {
	if len(nounRatings) == 0 || len(adjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingsLists()

	sort.Sort(nouns)
	sort.Sort(adjectives)

	nouns = splitNegative(nouns)
	adjectives = splitNegative(adjectives)

	if len(nouns) == 0 || len(adjectives) == 0 {
		return "No bad insults"
	}

	adj := adjectives[rand.Intn(len(adjectives))].Key
	noun := nouns[rand.Intn(len(nouns))].Key

	return createInsult(channel, target, adj, noun)
}

func BestInsult(channel string, target string) string {
	if len(nounRatings) == 0 || len(adjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingsLists()

	sort.Sort(sort.Reverse(nouns))
	sort.Sort(sort.Reverse(adjectives))

	adj := adjectives[0].Key
	noun := nouns[0].Key

	return createInsult(channel, target, adj, noun)
}

func WorstInsult(channel string, target string) string {
	if len(nounRatings) == 0 || len(adjectiveRatings) == 0 {
		return "No rated insults"
	}

	adjectives, nouns := getRatingsLists()

	sort.Sort(nouns)
	sort.Sort(adjectives)

	adj := adjectives[0].Key
	noun := nouns[0].Key

	return createInsult(channel, target, adj, noun)
}

func getRatingsLists() (pairList, pairList) {
	var nouns pairList
	for k, v := range nounRatings {
		nouns = append(nouns, pair{k, v})
	}

	var adjectives pairList
	for k, v := range adjectiveRatings {
		adjectives = append(adjectives, pair{k, v})
	}
	sort.Sort(sort.Reverse(adjectives))

	return adjectives, nouns
}

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
