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
	fmt.Println("handleGetReviews called")
}

func handlePostReviews(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlePostReviews called")
}
