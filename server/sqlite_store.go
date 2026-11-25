package main

import (
	"database/sql"

	_ "modernc.org/sqlite" // registers the "sqlite" driver with database/sql
)

// SQLiteStore stores reviews in a SQLite database file so they persist between runs.
type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			servery TEXT,
			dish TEXT,
			rating INTEGER,
			comment TEXT
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) GetReviews(servery string) ([]Review, error) {
	// Use ? placeholders instead of string concatenation to prevent SQL injection.
	rows, err := s.db.Query("SELECT servery, dish, rating, comment FROM reviews WHERE servery = ?", servery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// we need to convert from SQL types (rows of TEXT and INTEGER columns) to Go structs (Review)
	matching := []Review{}
	for rows.Next() {
		var r Review
		if err := rows.Scan(&r.Servery, &r.Dish, &r.Rating, &r.Comment); err != nil {
			return nil, err
		}
		matching = append(matching, r)
	}
	return matching, rows.Err()
}

func (s *SQLiteStore) AddReview(review Review) error {
	_, err := s.db.Exec(
		"INSERT INTO reviews (servery, dish, rating, comment) VALUES (?, ?, ?, ?)",
		review.Servery, review.Dish, review.Rating, review.Comment,
	)
	return err
}
