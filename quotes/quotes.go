package quotes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"api.deveshanand.com/db"
	types "api.deveshanand.com/quotes/types"

	"github.com/gin-gonic/gin"
)

func QuotesHandler(c *gin.Context) {
	//connect db
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	author := c.Query("author")
	date := c.Query("date")

	var quotes *sql.Rows
	var err error
	if author != "" && date != "" {
		quotes, err = con.Query("select * from quotes where author=? and date=?", author, date)
	} else if author != "" {
		fmt.Println("author is not empty")
		quotes, err = con.Query("select * from quotes where author LIKE ? AND active=1", author)
	} else if date != "" {
		quotes, err = con.Query("select * from quotes where date=? AND active=1", date)
	} else {
		quotes, err = con.Query("select * from quotes where active=1")
	}

	if err != nil {
		panic(err)
	}
	defer quotes.Close()

	var qts []types.Quote
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
		q := types.Quote{Id: id, Quote: quote, Author: author, Sub_by: sub_by, Date: date, Active: active}

		qts = append(qts, q)
	}
	c.JSONP(200, qts)
}

func SubmitQuote(c *gin.Context) {
	var userQuote types.PostData
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
				"message": "Quote already exists",
			})
			return
		} else {
			// no row found
			AddQuote(userQuote, 0)
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

func AddQuote(quote types.PostData, active int) {
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	quotes, err := con.Query("insert into quotes (quote, author, sub_by, active) values (?, ?, ?, ?)", quote.Quote, quote.Author, quote.Sub_by, active)
	if err != nil {
		panic(err)
	}
	defer quotes.Close()
}

func GetAllAuthors(c *gin.Context) {
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	authors, err := con.Query("select distinct author from quotes where active=1")
	if err != nil {
		panic(err)
	}
	defer authors.Close()

	var auths []string
	for authors.Next() {
		var author string
		_ = authors.Scan(&author)
		auths = append(auths, author)
	}
	c.JSONP(200, auths)
}
