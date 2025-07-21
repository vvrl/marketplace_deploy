package storage

import "database/sql"

type Storage struct {
	DB   *sql.DB
	User UserStorage
	Ad   AdStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		DB:   db,
		User: NewUserStorage(db),
		Ad:   NewAdStorage(db),
	}
}
