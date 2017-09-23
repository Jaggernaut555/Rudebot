package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const (
	wordDirectory = "words"
	nounPath      = "words/nouns/all.txt"
	adjectivePath = "words/adjectives/all.txt"
	adverbPath    = "words/adverbs/all.txt"
	verbPath      = "words/verbs/all.txt"
)

var (
	NumAdjectives int
	Adjectives    []string
	NumAdverbs    int
	Adverbs       []string
	NumNouns      int
	Nouns         []string
	NumVerbs      int
	Verbs         []string
)

func InitInsults() {
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

func NewInsult(target string) string {
	insult := target + " is a"
	adj := Adjectives[RandomInt(NumAdjectives)]
	noun := Nouns[RandomInt(NumNouns)]

	if StartsWithVowel(adj) {
		insult += "n"
	}
	insult += " " + adj
	insult += " " + noun

	return insult
}

func RandomInt(max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max)
}
