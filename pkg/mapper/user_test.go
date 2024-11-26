package mapper

import (
	"testing"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/stretchr/testify/require"
)

func TestMapUserAuthInputToEntity(t *testing.T) {
	testCases := map[string]struct {
		input    *request.UserAuthInput
		expected *entity.User
	}{
		"valid_input": {
			input: &request.UserAuthInput{
				Username: "testuser",
				Password: "password123",
			},
			expected: &entity.User{
				Username: "testuser",
				Password: "password123",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := MapUserAuthInputToEntity(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
