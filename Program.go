package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./vue", true)))

	router.GET("/suggest/:word", suggestionHandler)
	router.PUT("/suggest/:add", permutationHandler)

	router.Run(":8081")
}

func suggestionHandler(context *gin.Context) {
	word := context.Param("word")
	corrected := SuggestCorrection(&word)
	fmt.Println(corrected)
	var suggestions []string
	if corrected != "" {
		suggestions = SuggestedWords(&corrected)
	} else {
		suggestions = SuggestedWords(&word)
	}
	if len(suggestions) > 5 {
		suggestions = suggestions[0:5]
	}
	context.JSON(http.StatusOK, suggestions)
}

func permutationHandler(context *gin.Context) {
	word := context.Param("add")
	splt := strings.Split(word, ";")
	freq := int64(50000000)
	AddCompoundToTrie(splt, &freq)
}
