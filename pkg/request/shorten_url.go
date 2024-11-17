package request

import (
	"fmt"
	"net/url"
)

type ShortenUrlRequest struct {
	Url         string `json:"url"`
	CustomShort string `json:"custom_short"`
	Expiration  string `json:"expiration"`
}

func (r *ShortenUrlRequest) Validate() error {
	if r.Url == "" {
		return fmt.Errorf("url is required")
	}
	_, err := url.ParseRequestURI(r.Url)
	if err != nil {
		return fmt.Errorf("invalid url format")
	}
	return nil
}
