package main

import (
	"bufio"
	"hash/fnv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var node *trieNode
var currentNode *trieNode

var generatedFrequencies map[string]int64
var generatedEdits map[uint32][]string

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

	generatedFrequencies = make(map[string]int64)
	generatedEdits = make(map[uint32][]string)
	loadFrequencies()
	println("Finished loading frequencies and edits")
	for word, frequency := range generatedFrequencies {
		add(&word, &word, &frequency)
	}
	println("Finished loading trie")
}

//Add function adds words + character entry to and Trie
func add(originalWord *string, word *string, frequency *int64) {
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
		add(originalWord, &safeRecursiveString, frequency)
	}
	currentNode = node
}

//SuggestedWords returns autocomplete array of three top words by frequency
func SuggestedWords(word *string) []SortPair {
	char := []rune(*word)
	for index := 0; index < len(char); index++ {
		nodeWord := currentNode.children[string(char[index])]
		for key := range nodeWord.words {
			println(key)
		}
		currentNode = &nodeWord
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
	//TODO: add properly working sort
	sort.Slice(set[:], func(i, j int) bool {
		return set[i].key > set[j].key
	})
	currentNode = node
	return set[0:3]
}

//edits form all deletes for a word
func edits(word *string, distance int, set map[string]string) map[string]string {
	distance++
	for index := 0; index < len(*word); index++ {
		char := []rune(*word)
		edited := string(append(char[:index], char[index+1:]...))
		if _, exists := set[edited]; !exists && edited != "" {
			set[edited] = edited
			if distance < 2 {
				edits(&edited, distance, set)
			}
		}
	}
	return set
}

func loadFrequencies() {
	file, err := os.Open("frequency_dictionary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	k := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		frequency, err := strconv.ParseInt(values[1], 10, 64)
		if err != nil {
			println(values[0])
			continue
		}
		values[0] = strings.TrimSuffix(values[0], "...+2 more")
		generatedFrequencies[values[0]] = frequency
		if k != 0 {
			edited := edits(&values[0], 0, make(map[string]string))
			for edit := range edited {
				editHash := hash(&edit)
				if existingEdit, exists := generatedEdits[editHash]; exists {
					existingEdit = append(existingEdit, values[0])
				} else {
					generatedEdits[editHash] = []string{values[0]}
				}
			}
		}
		k++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func hash(word *string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(*word))
	return h.Sum32()
}
