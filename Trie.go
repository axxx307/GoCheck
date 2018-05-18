package main

var node *trieNode
var currentNode *trieNode

//trieNode is a node in Trie structure
type trieNode struct {
	children map[string]trieNode
	symbol   string
	words    map[string]int64
}

type sortPair struct {
	key   string
	value int64
}

//Init construtor
func Init() {
	node = &trieNode{children: make(map[string]trieNode), symbol: "", words: make(map[string]int64)}
	currentNode = node
}

//Add function adds words + character entry to and Trie
func Add(word *string, frequency *int64) {
	char := []rune(*word)
	safeSubstring := char[0:1]
	if currentNode.children[*word].symbol == "" {
		newNode := &trieNode{children: make(map[string]trieNode), symbol: string(safeSubstring), words: make(map[string]int64)}
		if len(*word) == 1 {
			newNode.words[*word] = *frequency
		}
		currentNode.children[string(safeSubstring)] = *newNode
	}
	if len(*word) > 1 {
		child := currentNode.children[string(safeSubstring)]
		currentNode = &child
		safeRecursiveString := string(char[0 : len(*word)-1])
		Add(&safeRecursiveString, frequency)
	}
	currentNode = node
}

//Get returns autocomplete array of three top words by frequency
func Get(word *string) [3]string {
	char := []rune(*word)
	for index := 0; index < len(char); index++ {
		nodeWord := currentNode.children[string(char[index])]
		println(nodeWord.symbol)
		currentNode = &nodeWord
	}
	//var sort []sortPair
	set := make(map[string]map[string]int64)
	for index, child := range currentNode.children {
		set[index] = child.words
	}
	for child := range set {
		println(child)
	}
	return [3]string{"", "", ""}
}
