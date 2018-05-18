package main

import (
	"fmt"
)

func main() {
	fmt.Println("sup")
	Init()
	st := "test"
	it := int64(15)
	Add(&st, &it)
	t := "tes"
	Get(&t)
}
