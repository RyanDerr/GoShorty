package mapper

import "github.com/RyanDerr/GoShorty/pkg/response"

func MapSignedJwtToResponse(signedJwt string) *response.JwtReponse {
	return &response.JwtReponse{
		Jwt: signedJwt,
	}
}
