package entity

import "time"

type ShortenUrl struct {
	BaseUrl    string        `json:"url"`
	Short      string        `json:"short"`
	Expiration time.Duration `json:"expiration"`
}

