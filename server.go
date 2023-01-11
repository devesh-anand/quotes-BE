package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/gin-gonic/gin"
)

func main() {
	// r := gin.Default()
	//routes
	http.HandleFunc("/", helloHandler)
	fmt.Printf("server on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello World!")
}
