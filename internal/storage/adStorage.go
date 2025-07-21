package storage

import (
	"context"
	"database/sql"
	"fmt"
	"marketplace/internal/logger"
	"marketplace/internal/models"
	"strings"
	"time"
)

type (
	AdStorage interface {
		CreateAdvertisement(ctx context.Context, ad *models.Advertisement) (*models.Advertisement, error)
		GetAdList(ctx context.Context, params models.ForListAdsParams) ([]*models.Advertisement, error)
	}

	adStorage struct {
		db *sql.DB
	}
)

func NewAdStorage(db *sql.DB) AdStorage {
	return &adStorage{db: db}
}

func (s *adStorage) CreateAdvertisement(ctx context.Context, ad *models.Advertisement) (*models.Advertisement, error) {
	query := `
	INSERT INTO ads (title, text, image_url, price, user_id)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctxWithTimeout, query, ad.Title, ad.Text, ad.ImageURL, ad.Price, ad.AuthorID).Scan(&ad.ID, &ad.CreatedAt)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (s *adStorage) GetAdList(ctx context.Context, params models.ForListAdsParams) ([]*models.Advertisement, error) {

	orderDirection := "ASC"
	if strings.ToUpper(params.Direction) == "DESC" {
		orderDirection = "DESC"
	}

	orderBy := "a.created_at"
	if params.Order == "a.price" {
		orderBy = "a.price"
	}

	query := fmt.Sprintf(`
	SELECT a.id, a.title, a.text, a.image_url, a.price, u.id AS author_id
	FROM ads a
	JOIN users u ON a.user_id = u.id
	WHERE a.price BETWEEN $1 AND $2
	ORDER BY %s %s
	LIMIT $3 OFFSET $4
	`, orderBy, orderDirection)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctxWithTimeout, query,
		params.MinPrice,
		params.MaxPrice,
		params.Limit,
		(params.Page-1)*params.Limit,
	)

	if err != nil {
		logger.Logger.Errorf("get ad list error: %v", err)
		return nil, err
	}

	ads := make([]*models.Advertisement, 0)

	for rows.Next() {
		temp := &models.Advertisement{}
		err = rows.Scan(
			&temp.ID,
			&temp.Title,
			&temp.Text,
			&temp.ImageURL,
			&temp.Price,
			&temp.AuthorID,
		)
		if err != nil {
			return nil, err
		}
		ads = append(ads, temp)
	}
	fmt.Println("storage")
	fmt.Println(ads)
	return ads, nil
}
