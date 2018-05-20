package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	router := gin.Default()

	router.LoadHTMLGlob("vue/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/suggest/:word", suggestionHandler)

	router.Run(":8080")
}

func suggestionHandler(context *gin.Context) {
	word := context.Param("word")
	suggestions := SuggestedWords(&word)
	context.JSON(http.StatusOK, suggestions)
}
