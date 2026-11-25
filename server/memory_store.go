package main

import "sync"

// MemoryStore stores reviews in a slice in memory.
// These reviews will only live for the life of the Go program, not between runs.
type MemoryStore struct {
	reviews []Review
	mu      sync.Mutex // protects reviews against concurrent reads/writes
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		reviews: []Review{},
	}
}

func (s *MemoryStore) GetReviews(servery string) ([]Review, error) {
	// we must prevent concurrent writers from changing the reviews while we try to read from it
	s.mu.Lock()         // lock before reading from our stored reviews
	defer s.mu.Unlock() // we'll unlock at the end of this function

	matching := []Review{}
	for _, review := range s.reviews {
		if review.Servery == servery {
			matching = append(matching, review)
		}
	}
	return matching, nil
}

func (s *MemoryStore) AddReview(review Review) error {
	// we must prevent concurrent writers from changing the stored list of reviews while we try to write to it
	// and we must prevent concurrent readers from getting an incorrect view of the data while we append a review
	s.mu.Lock()         // lock before writing to our stored reviews
	defer s.mu.Unlock() // we'll unlock at the end of this function

	s.reviews = append(s.reviews, review)
	return nil
}
