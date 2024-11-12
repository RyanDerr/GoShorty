package response

import "fmt"

type ShortenUrlResponse struct {
	Url         string `json:"url"`
	CustomShort string `json:"custom_short"`
	Expiration  string `json:"expiration"`
}

func (r *ShortenUrlResponse) String() string {
	str := fmt.Sprintf("URL: %s\nCustomShort: %s\nExpiration: %s", r.Url, r.CustomShort, r.Expiration)
	return str
}
