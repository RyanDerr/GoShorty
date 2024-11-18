package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAuthInput_Validate(t *testing.T) {
	testCases := map[string]struct {
		input           *UserAuthInput
		wantErr         bool
		wantErrContains string
	}{
		"valid_input": {
			input: &UserAuthInput{
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		"missing_username": {
			input: &UserAuthInput{
				Password: "password123",
			},
			wantErr:         true,
			wantErrContains: "username is required",
		},
		"missing_password": {
			input: &UserAuthInput{
				Username: "testuser",
			},
			wantErr:         true,
			wantErrContains: "password is required",
		},
		"missing_username_and_password": {
			input:           &UserAuthInput{},
			wantErr:         true,
			wantErrContains: "username is required",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrContains)
				return
			}

			require.NoError(t, err)
		})
	}
}
