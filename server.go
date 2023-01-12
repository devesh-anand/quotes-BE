//go:build !production && production
// +build !production,production

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"quotes-BE/db"
	"quotes-BE/routes"
)

func main() {
	godotenv.Load()

	//connect db
	con, conerr := db.GetConnection()
	if conerr != nil {
		panic(conerr)
	}
	defer con.Close()

	//router and endpoints
	r := routes.SetupRouter()
	r.Run(os.Getenv("PORT"))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
