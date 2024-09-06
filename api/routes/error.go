package routes

type ErrorResponse struct {
	Error string `json:"error"`
}

type RateLimitExceededResponse struct {
	Error          string `json:"error"`
	RateLimitReset string `json:"rate_limit_reset"`
}
