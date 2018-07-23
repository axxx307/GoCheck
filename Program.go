package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	text := "handlin"
	fmt.Println(SuggestCorrection(&text))

	// router := gin.Default()

	// // Serve frontend static files
	// router.Use(static.Serve("/", static.LocalFile("./vue", true)))

	// router.GET("/suggest/:word", suggestionHandler)

	// router.Run(":8081")
}

func suggestionHandler(context *gin.Context) {
	word := context.Param("word")
	suggestions := SuggestedWords(&word)
	if len(suggestions) > 5 {
		suggestions = suggestions[0:5]
	}
	context.JSON(http.StatusOK, suggestions)
}

///при генерации документа записываем его в монгу
///вместе с доступными пермутациями
///из монги составить префиксное дерево на основе термов и держать их в памяти/кэше
///при запросах, сначала доставать по пермутациям вместе с носновной информацией,
///в дереве находить СЛЕДУЮЩЕЕ слово, базируясь на метадате найденного слова
///Handling Catalog, Part
///Handling Catalog rank 2, Part rank 1

///записали хэши делитедов в монгу
///записали данные
///нашли слово
