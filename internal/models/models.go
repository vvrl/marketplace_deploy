package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	HashPassword string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Advertisement struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url"`
	Price     float64   `json:"price"`
	AuthorID  int       `json:"author_id"`
	IsMine    bool      `json:"is_mine"`
	CreatedAt time.Time `json:"created_at"`
}

type ForListAdsParams struct {
	MinPrice  float64
	MaxPrice  float64
	Order     string
	Direction string
	Page      int
	Limit     int
}
