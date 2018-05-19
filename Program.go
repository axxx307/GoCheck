package main

func main() {
	Init()
	word := "tes"
	completed := SuggestedWords(&word)
	for _, item := range completed {
		println(item.key)
	}
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print("Enter text: ")
	// 	text, _ := reader.ReadString('\n')
	// 	completed := Get(&text)
	// 	for item := range completed {
	// 		print(item)
	// 	}
	// }
}
