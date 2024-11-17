package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testCases := map[string]struct {
		input           *ShortenUrlRequest
		wantErr         bool
		wantErrContains string
	}{
		"valid_request": {
			input: &ShortenUrlRequest{
				Url:         "http://example.com",
				CustomShort: "exmpl",
				Expiration:  "1h",
			},
			wantErr: false,
		},
		"valid_no_custom_short": {
			input: &ShortenUrlRequest{
				Url:        "http://example.com",
				Expiration: "1h",
			},
			wantErr: false,
		},
		"valid_no_expiration": {
			input: &ShortenUrlRequest{
				Url:         "http://example.com",
				CustomShort: "exmpl",
			},
			wantErr: false,
		},
		"invalid_url": {
			input: &ShortenUrlRequest{
				Url: "invalid_url",
			},
			wantErr:         true,
			wantErrContains: "invalid url format",
		},
		"missing_url": {
			input: &ShortenUrlRequest{},
			wantErr:         true,
			wantErrContains: "url is required",
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
