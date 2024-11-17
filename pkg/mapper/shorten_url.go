package mapper

import (
	"fmt"
	"time"

	urlEntity "github.com/RyanDerr/GoShorty/internal/domain/entity/url"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/google/uuid"
)

const (
	// defaultExpiration is the default expiration time for a shortened URL if none is provided
	defaultExpiration = time.Hour * 24
)

func MapShortenUrlRequestToEntity(req *request.ShortenUrlRequest) (*urlEntity.ShortenUrl, error) {
	var cShort string
	if req.CustomShort == "" {
		cShort = uuid.New().String()[:6]
	} else {
		cShort = req.CustomShort
	}

	var exp time.Duration
	if req.Expiration != "" {
		expiration, err := time.ParseDuration(req.Expiration)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expiration duration: %v", err)
		}
		exp = expiration
	} else {
		exp = defaultExpiration
	}

	return &urlEntity.ShortenUrl{
		BaseUrl:    req.Url,
		Short:      cShort,
		Expiration: exp,
	}, nil
}

func MapShortenUrlEntityToResponse(entity *urlEntity.ShortenUrl) *response.ShortenUrlResponse {
	return &response.ShortenUrlResponse{
		Url:         entity.BaseUrl,
		CustomShort: entity.Short,
		Expiration:  entity.Expiration.String(),
	}
}