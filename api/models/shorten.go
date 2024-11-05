package models

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