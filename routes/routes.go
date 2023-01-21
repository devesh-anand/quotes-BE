package routes

import (
	"fmt"

	"quotes-BE/quotes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())
	//routes
	r.GET("/", helloHandler)
	fmt.Printf("server on port 8080")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Quotes routes
	r.GET("/quotes", quotes.QuotesHandler)
	r.POST("/quote", quotes.SubmitQuote)

	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(202, gin.H{
		"data": "Hello World",
	})

	fmt.Print("Hello World!")
}
