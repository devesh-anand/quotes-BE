package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"api.deveshanand.com/db"
	"api.deveshanand.com/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	//cron to add a quote every 6 days so as to prevent db from sleeping
	// not needed as we are using turso db
	// cron.SubmitCron()

	//connect db
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	//router and endpoints
	r := routes.SetupRouter()
	
	log.Printf("Server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
