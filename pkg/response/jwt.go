package response

type JwtReponse struct {
	Jwt string
}

func (r *JwtReponse) String() string {
	return r.Jwt
}
