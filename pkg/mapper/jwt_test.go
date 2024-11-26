package mapper

import (
    "testing"

    "github.com/RyanDerr/GoShorty/pkg/response"
    "github.com/stretchr/testify/require"
)

const (
	dummyJwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

func TestMapSignedJwtToResponse(t *testing.T) {
    testCases := map[string]struct {
        input    string
        expected *response.JwtResponse
    }{
        "valid_jwt": {
            input: dummyJwtToken,
            expected: &response.JwtResponse{
                Jwt: dummyJwtToken,
            },
        },
    }

    for name, tc := range testCases {
        t.Run(name, func(t *testing.T) {
            result := MapSignedJwtToResponse(tc.input)
            require.Equal(t, tc.expected, result)
        })
    }
}