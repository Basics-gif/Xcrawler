package storage

import (
	"database/sql"
	"fmt"

	"github.com/Basics-gif/go_crawler/internal/browser"
	_ "modernc.org/sqlite"
)

const (
	PlanPaid = "paid"
	PlanFree = "free"

	durationThreshold = 20 * 60
	viewsThreshold    = 100_000
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("could open database: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS videos (
		page_url TEXT PRIMARY KEY,
		title 	 TEXT NOT NULL,
		duration INTEGER,
		plan 		 TEXT NOT NULL CHECK(plan IN ('paid', 'free')),
		created_at DETETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not create schema: %v", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func DeterminePlan(duration int, views int64) string {
	if duration > durationThreshold || views > viewsThreshold {
		return PlanPaid
	}
	return PlanFree
}

func (s *Store) Save(v browser.Video) error {
	plan := DeterminePlan(v.Duration, v.Views)

	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO videos (page_url, title, duration, plan) VALUES (?, ?, ?, ?)`,
		v.PageURL, v.Title, v.Duration, plan,
	)
	if err != nil {
		return fmt.Errorf("could not save video %s: %v", v.PageURL, err)
	}
	return nil

}

func (s *Store) SaveAll(list *browser.VideoList) error {
	for _, v := range list.Videos {
		if err := s.Save(v); err != nil {
			return err
		}
	}
	return nil
}
