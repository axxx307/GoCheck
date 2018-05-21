package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./vue", true)))

	router.GET("/suggest/:word", suggestionHandler)

	router.Run(":8081")
}

func suggestionHandler(context *gin.Context) {
	word := context.Param("word")
	suggestions := SuggestedWords(&word)
	context.JSON(http.StatusOK, suggestions)
}
