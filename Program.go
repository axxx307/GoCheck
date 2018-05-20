package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	router := gin.Default()

	router.GET("/suggest/:word", suggestionHandler)

	router.Run(":8080")
}

func suggestionHandler(context *gin.Context) {
	word := context.Param("word")
	suggestions := SuggestedWords(&word)
	context.String(http.StatusOK, strings.Join(suggestions, " "))
}
