package main

import (
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

var cache = map[string]map[string]int{}
var cacheRemoveKeys = map[string]bool{}

func GetTopWords(numberOfWords int, sessionID string, text string) ([]TopWord, string) {
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	delete(cacheRemoveKeys, sessionID)
	words := cache[sessionID]
	if words == nil {
		words = map[string]int{}
	}
	// make it lowercase, so we don't have to worry about uppercase and lowercase letters'
	text = strings.ToLower(text)
	// remove enter (new line)
	stopWords := []string{"\n", "\r", "\t", "(", ")", "\"", "'"}

	for _, stopWord := range stopWords {
		text = strings.ReplaceAll(text, stopWord, " ")
	}
	tempWords := strings.Split(text, " ")
	for _, word := range tempWords {
		// remove punctuation and space
		word = strings.Trim(word, " ?.!:;,()[]{}\t")
		word = strings.ReplaceAll(word, " ", "")
		if word == "" {
			continue
		}
		_, exist := words[word]
		if exist {
			words[word]++
		} else {
			words[word] = 1
		}
	}
	mapKeys := []string{}
	for key := range words {
		mapKeys = append(mapKeys, key)
	}
	sort.Slice(mapKeys, func(i, j int) bool {
		return words[mapKeys[i]] > words[mapKeys[j]]
	})
	topWords := []TopWord{}
	for i := 0; i < numberOfWords && i < len(mapKeys); i++ {
		topWord := TopWord{
			Word:  mapKeys[i],
			Count: words[mapKeys[i]],
		}
		topWords = append(topWords, topWord)
	}
	cache[sessionID] = words
	go removeCache(sessionID)
	return topWords, sessionID
}

type TopWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

func removeCache(sessionID string) {
	cacheRemoveKeys[sessionID] = true
	time.Sleep(time.Minute * 5)
	for k := range cacheRemoveKeys {
		if k == sessionID {
			delete(cache, k)
			return
		}
	}
}
