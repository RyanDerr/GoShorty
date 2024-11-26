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
		"short_username": {
			input: &UserAuthInput{
				Username: "usr",
				Password: "password123",
			},
			wantErr:         true,
			wantErrContains: "username must be at least 4 characters",
		},
		"long_username": {
			input: &UserAuthInput{
				Username: "thisisaverylongusernamethatexceedsthemaxlimit",
				Password: "password123",
			},
			wantErr:         true,
			wantErrContains: "username must be at most 32 characters",
		},
		"invalid_username": {
			input: &UserAuthInput{
				Username: "invalid_user!",
				Password: "password123",
			},
			wantErr:         true,
			wantErrContains: "username must be alphanumeric",
		},
		"short_password": {
			input: &UserAuthInput{
				Username: "testuser",
				Password: "short",
			},
			wantErr:         true,
			wantErrContains: "password must be at least 8 characters",
		},
		"long_password": {
			input: &UserAuthInput{
				Username: "testuser",
				Password: "thisisaverylongpasswordthatexceedsthemaxlimit",
			},
			wantErr:         true,
			wantErrContains: "password must be at most 32 characters",
		},
		"invalid_password": {
			input: &UserAuthInput{
				Username: "testuser",
				Password: "invalid password!",
			},
			wantErr:         true,
			wantErrContains: "password contains invalid characters",
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
