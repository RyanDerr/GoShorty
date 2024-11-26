package response

type JwtResponse struct {
	Jwt string
}

func (r *JwtResponse) String() string {
	return r.Jwt
}
