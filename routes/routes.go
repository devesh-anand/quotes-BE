package routes

import (
	"fmt"

	"quotes-BE/quotes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"https://quotes.deveshanand.com", "http://quotes.deveshanand.com"},
	// 	// AllowOrigins:     []string{"http://localhost:3000", "https://quotes.deveshanand.com"},
	// 	AllowMethods:     []string{"OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "X-Requested-With", "Authorization", "Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

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
