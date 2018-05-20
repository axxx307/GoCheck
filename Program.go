package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Init()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		println(text)
		completed := SuggestedWords(&text)
		for index := 0; index < 2; index++ {
			println(completed[index].value)
		}
	}
}
