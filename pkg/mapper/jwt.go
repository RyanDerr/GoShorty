package mapper

import "github.com/RyanDerr/GoShorty/pkg/response"

func MapSignedJwtToResponse(signedJwt string) *response.JwtResponse {
	return &response.JwtResponse{
		Jwt: signedJwt,
	}
}
