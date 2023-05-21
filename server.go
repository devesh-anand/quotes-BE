package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"api.deveshanand.com/db"
	"api.deveshanand.com/quotes/cron"
	"api.deveshanand.com/routes"
)

func main() {
	godotenv.Load()

	//cron to add a quote every 6 days so as to prevent db from sleeping
	cron.SubmitCron()

	//connect db
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	//router and endpoints
	r := routes.SetupRouter()
	r.Run(":" + os.Getenv("PORT"))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
