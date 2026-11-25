package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// Must match the servery ids hardcoded in page.tsx.
var serveryIDs = []string{"baker", "north", "seibel", "south", "south-main", "west"}

// Review mirrors the JavaScript object on line 43 of page.tsx.
// The tags on this struct (`json:"servery"`) are used when the json package converts this Go struct to/from a JSON object
type Review struct {
	Servery string `json:"servery"`
	Dish    string `json:"dish"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// ReviewStore defines how reviews are saved and loaded.
// Using an interface lets our HTTP handlers work with any storage backend.
type ReviewStore interface {
	GetReviews(servery string) ([]Review, error)
	AddReview(review Review) error
}

// Server holds our storage dependency and handles HTTP requests.
type Server struct {
	store ReviewStore
}

func main() {
	fmt.Println("servery-reviews backend running...")

	// we can change between a MemoryStore and SQLiteStore seamlessly since they both fulfill the ReviewStore interface
	store, err := NewSQLiteStore("reviews.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer store.Close()

	srv := &Server{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/reviews", srv.handleGetReviews)
	mux.HandleFunc("POST /api/reviews", srv.handlePostReviews)

	// I chose 8080 because the README says backend requests are sent to this port.
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) handleGetReviews(w http.ResponseWriter, r *http.Request) {
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

	matching, err := s.store.GetReviews(servery)
	if err != nil {
		http.Error(w, "Failed to load reviews.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// notice that the response is a list of Reviews, which matches the expected return type of getReviews in page.tsx
	// it's bad practice, but we'll ignore the error from Encode for now
	json.NewEncoder(w).Encode(matching)
}

func (s *Server) handlePostReviews(w http.ResponseWriter, r *http.Request) {
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

	// validate and clean request data
	review.Dish = strings.TrimSpace(review.Dish)
	review.Comment = strings.TrimSpace(review.Comment)

	if !isServery(review.Servery) {
		http.Error(w, "Unknown servery.", http.StatusBadRequest)
		return
	}
	if review.Dish == "" {
		http.Error(w, "Dish name is required.", http.StatusBadRequest)
		return
	}
	if review.Rating < 1 || review.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5.", http.StatusBadRequest)
		return
	}

	if err := s.store.AddReview(review); err != nil {
		http.Error(w, "Failed to save review.", http.StatusInternalServerError)
		return
	}

	// we don't need to return any data for this POST request, just provide the client a success code
	w.WriteHeader(http.StatusCreated)
}

func isServery(id string) bool {
	return slices.Contains(serveryIDs, id)
}
