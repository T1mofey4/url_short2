package app

import "time"

type URL struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	LongURL   string    `json:"long_url"`
	TTL       time.Time `json:"ttl"`
	CreatedAt time.Time `json:"created_at"`
}
