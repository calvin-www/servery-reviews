package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("servery-reviews backend running...")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/reviews", handleGetReviews)
	mux.HandleFunc("POST /api/reviews", handlePostReviews)

	// I chose 8080 because the README says backend requests are sent to this port.
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

func handleGetReviews(w http.ResponseWriter, r *http.Request) {
	// this function should execute when you run the following command in the terminal (while running this main.go file):
	// curl -X GET --location "http://localhost:8080/api/reviews?servery=north"
	fmt.Println("handleGetReviews called")
}

func handlePostReviews(w http.ResponseWriter, r *http.Request) {
	// this function should execute when you run the following command in the terminal (while running this main.go file):
	// curl -X POST --json '{ "servery": "north", "dish": "Pad Thai", "rating": 4, "comment": "Pretty good" }' --location "http://localhost:8080/api/reviews"
	fmt.Println("handlePostReviews called")
}
