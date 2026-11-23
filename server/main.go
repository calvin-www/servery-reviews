package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Review mirrors the JavaScript object on line 43 of page.tsx.
// The tags on this struct (`json:"servery"`) are used when the json package converts this Go struct to/from a JSON object
type Review struct {
	Servery string `json:"servery"`
	Dish    string `json:"dish"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// declare a global variable to hold all reviews
// these reviews will only live for the life of the go program, not between runs
var reviews = []Review{}

func main() {
	fmt.Println("servery-reviews backend running...")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/reviews", handleGetReviews)
	mux.HandleFunc("POST /api/reviews", handlePostReviews)

	// I chose 8080 because the README says backend requests are sent to this port.
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleGetReviews(w http.ResponseWriter, r *http.Request) {
	// this function should execute when you run the following command in the terminal (while running this main.go file):
	// curl -X GET --location "http://localhost:8080/api/reviews?servery=north"
	fmt.Println("handleGetReviews called")

	// need to read "?servery=north" query parameter from REST API request
	// searched up "golang read http request params"
	servery := r.URL.Query().Get("servery")
	fmt.Println(servery)

	if servery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	matching := []Review{}
	for _, review := range reviews {
		if review.Servery == servery {
			matching = append(matching, review)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// notice that the response is a list of Reviews, which matches the expected return type of getReviews in page.tsx
	// it's bad practice, but we'll ignore the error from Encode for now
	json.NewEncoder(w).Encode(matching)
}

func handlePostReviews(w http.ResponseWriter, r *http.Request) {
	// this function should execute when you run the following command in the terminal (while running this main.go file):
	// curl -X POST --json '{ "servery": "north", "dish": "Pad Thai", "rating": 4, "comment": "Pretty good" }' --location "http://localhost:8080/api/reviews"
	fmt.Println("handlePostReviews called")

	// need to read request body
	// searched up "golang read http request body"
	// also defined "Review" struct above
	defer r.Body.Close()

	var review Review
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&review); err != nil {
		fmt.Println(err)
	}

	// add review from request to our stored reviews
	reviews = append(reviews, review)

	// we don't need to return any data for this POST request, just provide the client a success code
	w.WriteHeader(http.StatusCreated)
}
