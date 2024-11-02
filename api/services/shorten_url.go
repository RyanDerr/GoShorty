package services

import (
	"fmt"
	"os"
	"time"

	"github.com/RyanDerr/GoShorty/api/models"
	"github.com/RyanDerr/GoShorty/api/urlparser"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func (s *URLService) ShortenURL(req *models.ShortenRequest) (*models.ShortenResponse, error) {
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}

	req.URL = urlparser.EnforceHTTPS(req.URL)
	id := generateID(req.CustomShort)

	if err := s.repo.CheckURLInUse(s.ctx, id); err != nil {
		return nil, err
	}

	expirationTime, err := parseExpiration(req.Expiration)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveURL(s.ctx, id, req.URL, expirationTime); err != nil {
		return nil, err
	}

	return populateResponse(req, id, expirationTime), nil
}

func validateURL(url string) error {
	if !govalidator.IsURL(url) {
		return fmt.Errorf("invalid url")
	}

	if !urlparser.RemoveDomainError(url) {
		return fmt.Errorf("url is already shortened")
	}

	return nil
}

func generateID(customShort string) string {
	if customShort == "" {
		return uuid.New().String()[:6]
	}
	return customShort
}

func parseExpiration(expiration string) (time.Duration, error) {
	if expiration == "" {
		return 24 * time.Hour, nil
	}
	return time.ParseDuration(expiration)
}

func populateResponse(req *models.ShortenRequest, id string, exp time.Duration) *models.ShortenResponse {
	return &models.ShortenResponse{
		URL:         req.URL,
		CustomShort: os.Getenv("DOMAIN") + "/" + id,
		Expiration:  exp.String(),
	}
}
