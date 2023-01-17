package routes

import (
	"fmt"
	"log"

	"quotes-BE/db"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Quote struct {
	Id     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Sub_by string `json:"sub_by"`
	Date   string `json:"date"`
	Active int    `json:"active"`
}

type PostData struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Sub_by string `json:"sub_by"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//routes
	r.GET("/", helloHandler)
	fmt.Printf("server on port 8080")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "po",
		})
	})

	r.GET("/str", func(c *gin.Context) {
		c.String(200, "string routes")
	})

	r.GET("/quotes", quotesHandler)

	r.POST("/submitquote", submitQuote)

	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(202, gin.H{
		"data": "Hello World",
	})

	fmt.Print("Hello World!")
}

func quotesHandler(c *gin.Context) {
	//connect db
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	quotes, err := con.Query("select * from quotes")
	if err != nil {
		panic(err)
	}
	defer quotes.Close()

	var qts []Quote
	for quotes.Next() {
		var (
			id     int
			quote  string
			author string
			sub_by string
			date   string
			active int
		)
		_ = quotes.Scan(&id, &quote, &author, &sub_by, &date, &active)
		q := Quote{Id: id, Quote: quote, Author: author, Sub_by: sub_by, Date: date, Active: active}

		qts = append(qts, q)
	}
	c.JSONP(200, qts)
}

func submitQuote(c *gin.Context) {
	var userQuote PostData
	if c.BindJSON(&userQuote) == nil {
		con, conerr := db.GetConnection()
		if conerr != nil {
			panic(conerr)
		}
		defer con.Close()

		check, err := con.Query("select * from quotes where quote=?", userQuote.Quote)
		if err != nil {
			log.Fatal(err)
		}

		if check.Next() {
			// row found
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Already Exists",
			})
			return
		} else {
			// no row found
			quotes, err := con.Query("insert into quotes (quote, author, sub_by, active) values (?, ?, ?, 0)", userQuote.Quote, userQuote.Author, userQuote.Sub_by)
			if err != nil {
				panic(err)
			}
			defer quotes.Close()
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   userQuote,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON format",
		})
		return
	}
}
