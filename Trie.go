package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var node *trieNode
var currentNode *trieNode
var words map[string]int
var deletions map[string][]string

//trieNode is a node in Trie structure
type trieNode struct {
	children       map[string]trieNode
	words          map[string]int64
	predictionType wordType
	nextWord       string
	metadata       string
	ranking        int32
}

type wordType struct {
	Concept  predictionType
	Subtitle predictionType
	Category predictionType
}

type predictionType string

//SortPair is used for sorting the resulting frequencies
//of all found autocomletion words
type SortPair struct {
	key   string
	value int64
}

//PairList A slice of pairs that implements sort.Interface to sort by values
type PairList []SortPair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].value < p[j].value }

//Init construtor
func Init() {
	node = &trieNode{children: make(map[string]trieNode), words: make(map[string]int64)}
	currentNode = node

	words = make(map[string]int)
	deletions = make(map[string][]string)

	loadFrequencies()
	println("Finished loading trie")

	println("Word training completed")
}

//addToTrie function adds words + character entry to and Trie
func addToTrie(originalWord *string, word *string, frequency *int64) {
	char := []rune(*word)
	safeSubstring := char[0:1]
	if existingNode, exists := currentNode.children[string(safeSubstring)]; exists {
		if _, wordExists := existingNode.words[*originalWord]; !wordExists && len(*word) == 1 {
			existingNode.words[*originalWord] = *frequency
		}
	} else {
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
		addToTrie(originalWord, &safeRecursiveString, frequency)
	}
	currentNode = node
}

//SuggestedWords returns autocomplete array of three top words by frequency
func SuggestedWords(word *string) []string {
	char := []rune(*word)
	for index := 0; index < len(char); index++ {
		if nodeWord, exists := currentNode.children[string(char[index])]; exists {
			currentNode = &nodeWord
		}
	}
	set := []SortPair{}
	for _, child := range currentNode.children {
		//get next childs words values
		for _, chld := range child.children {
			for word, frequency := range chld.words {
				set = append(set, SortPair{key: word, value: frequency})
			}
		}

		//get closest word and their frequencies
		for word, frequency := range child.words {
			set = append(set, SortPair{key: word, value: frequency})
		}
	}
	currentNode = node

	//Add sort to determent first five results
	sorted := make(PairList, len(set))
	for index, pair := range set {
		sorted[index] = pair
	}
	sort.Sort(sort.Reverse(sorted))

	keys := make([]string, 0, len(sorted))
	for _, k := range sorted {
		keys = append(keys, k.key)
	}
	return keys
}

func SuggestCorrection(word *string) string {
	if _, exists := words[*word]; exists {
		return *word
	}
	permutations := edits([]rune(*word), 2)
	for _, item := range permutations {
		h := sha1.New()
		h.Write([]byte(item))
		hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
		if items, exists := deletions[hash]; exists {
			freqs := make([]int, 0)
			for _, value := range items {
				freqs = append(freqs, words[value])
			}
			return items[0]
		}
	}
	return ""
}

//loads frequencies from txt file and adds them to trie
func loadFrequencies() {
	file, err := os.Open("frequency.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		frequency, err := strconv.ParseInt(values[1], 10, 64)
		if err != nil {
			println(err, values[0])
			continue
		}
		permutations := edits([]rune(values[0]), 0)
		for _, item := range permutations {
			println(item)
		}
		createArray(permutations, values[0])
		addToTrie(&values[0], &values[0], &frequency)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func createArray(permutations []string, word string) {
	for _, perm := range permutations {
		h := sha1.New()
		h.Write([]byte(perm))
		hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
		if _, exists := deletions[hash]; !exists {
			deletions[hash] = []string{word}
		} else {
			deletions[hash] = append(deletions[hash], word)
		}
	}
}

// only delete one, no transposes, replaces and inserts
func edits(q []rune, ed int) (v []string) {
	v = append(v, string(q))
	ed++

	for i := 0; i < len(q); i++ {
		x := remove(q, i)
		v = append(v, string(x))
		if ed < 2 {
			v = append(v, edits(x, ed)...)
		}
	}
	return
}

func remove(runes []rune, i int) []rune {
	var v = append([]rune{}, runes[:i]...)
	return append(v, runes[i+1:]...)
}
