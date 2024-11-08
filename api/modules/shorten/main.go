package shorten

import "fmt"

type ShortenRequest struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiration  string `json:"expiration"`
}

type ShortenResponse struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiration  string `json:"expiration"`
}

func (s *ShortenRequest) String() string {
	return fmt.Sprintf("URL: %s\nShort: %s\nExpiration: %s\n", s.URL, s.CustomShort, s.Expiration)
}

func (s *ShortenResponse) String() string {
	return fmt.Sprintf("URL: %s\nShort: %s\nExpiration: %s\n", s.URL, s.CustomShort, s.Expiration)
}
