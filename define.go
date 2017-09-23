package main

import (
	"fmt"
	"log"

	"github.com/karan/vocabulary"
	"github.com/kennygrant/sanitize"
)

const (
	BigHugeLabsApiKey = ""
	WordnikApiKey     = ""
)

func DefineWord(word string) string {
	c := &vocabulary.Config{BigHugeLabsApiKey: BigHugeLabsApiKey, WordnikApiKey: WordnikApiKey}

	v, err := vocabulary.New(c)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	definitions, err := v.Meanings(word)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	reply := "```\n" + word + "\n"

	for i, def := range definitions {
		reply += fmt.Sprintf("%d. %s\n", i+1, sanitize.HTML(def))
	}

	reply += "```"

	return reply
}
