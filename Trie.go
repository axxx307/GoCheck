package main

import (
	"sort"
)

var node *trieNode
var currentNode *trieNode

//trieNode is a node in Trie structure
type trieNode struct {
	children map[string]trieNode
	words    map[string]int64
}

//SortPair is used for sorting the resulting frequencies
//of all found autocomletion words
type SortPair struct {
	key   string
	value int64
}

//Init construtor
func Init() {
	node = &trieNode{children: make(map[string]trieNode), words: make(map[string]int64)}
	currentNode = node
}

//Add function adds words + character entry to and Trie
func Add(originalWord *string, word *string, frequency *int64) {
	char := []rune(*word)
	safeSubstring := char[0:1]
	if _, exists := currentNode.children[string(safeSubstring)]; !exists {
		newNode := &trieNode{children: make(map[string]trieNode), words: make(map[string]int64)}
		if len(*word) == 1 {
			newNode.words[*originalWord] = *frequency
		}
		currentNode.children[string(safeSubstring)] = *newNode
	}
	if len(*word) > 1 {
		child := currentNode.children[string(safeSubstring)]
		currentNode = &child
		safeRecursiveString := string(char[1:len(*word)])
		Add(originalWord, &safeRecursiveString, frequency)
	}
	currentNode = node
}

//Get returns autocomplete array of three top words by frequency
func Get(word *string) []SortPair {
	char := []rune(*word)
	for index := 0; index < len(char); index++ {
		nodeWord := currentNode.children[string(char[index])]
		currentNode = &nodeWord
	}
	set := []SortPair{}
	for _, child := range currentNode.children {
		for word, frequency := range child.words {
			set = append(set, SortPair{key: word, value: frequency})
		}
	}
	sort.Slice(set[:], func(i, j int) bool {
		return set[i].key > set[j].key
	})
	currentNode = node
	return set[0:2]
}
