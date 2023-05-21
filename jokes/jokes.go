package jokes

import (
	"log"
	"net/http"

	"api.deveshanand.com/db"

	"github.com/gin-gonic/gin"
)

// type Joke
type JokeData struct {
	Joke   string `json:"quote"`
	Author string `json:"author"`
	Sub_by string `json:"sub_by"`
}

func NewJoke(c *gin.Context) {
	var joke JokeData
	if c.BindJSON(&joke) == nil {
		con, conerr := db.GetConnection()
		if conerr != nil {
			panic(conerr)
		}
		defer con.Close()

		check, err := con.Query("select * from jokes where joke=?", joke.Joke)
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
			quotes, err := con.Query("insert into jokes (joke, author, sub_by, active) values (?, ?, ?, 0)", joke.Joke, joke.Author, joke.Sub_by)
			if err != nil {
				panic(err)
			}
			defer quotes.Close()
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   joke,
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
