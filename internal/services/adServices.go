package services

import (
	"context"
	"fmt"
	"marketplace/internal/logger"
	"marketplace/internal/models"
	"marketplace/internal/storage"
)

type (
	AdService interface {
		PostAd(ctx context.Context, title, text, imageURL string, price float64, userID int) (*models.Advertisement, error)
		GetAdList(ctx context.Context, params models.ForListAdsParams, userID int) ([]*models.Advertisement, error)
	}

	adService struct {
		storage storage.AdStorage
	}
)

func NewAdService(storage storage.AdStorage) AdService {
	return &adService{storage: storage}
}

func (s *adService) PostAd(ctx context.Context, title, text, imageURL string, price float64, userID int) (*models.Advertisement, error) {
	ad := &models.Advertisement{
		Title:    title,
		Text:     text,
		ImageURL: imageURL,
		Price:    price,
		AuthorID: userID,
		IsMine:   true,
	}

	ad, err := s.storage.CreateAdvertisement(ctx, ad)
	if err != nil {
		logger.Logger.Errorf("create advertisement error: %v", err)
		return nil, err
	}

	return ad, nil
}

func (s *adService) GetAdList(ctx context.Context, params models.ForListAdsParams, userID int) ([]*models.Advertisement, error) {
	ads, err := s.storage.GetAdList(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, ad := range ads {
		ad.IsMine = userID != 0 && ad.AuthorID == userID
	}
	fmt.Println("services")
	return ads, nil
}
