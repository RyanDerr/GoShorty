package mapper

import (
	"testing"
	"time"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/stretchr/testify/require"
)

func TestMapShortenUrlRequestToEntity(t *testing.T) {
	testCases := map[string]struct {
		input           *request.ShortenUrlRequest
		response        *entity.ShortenUrl
		wantErr         bool
		wantErrContains string
	}{
		"valid_request": {
			input: &request.ShortenUrlRequest{
				Url:         "http://example.com",
				CustomShort: "exmpl",
				Expiration:  "1h",
			},
			response: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			wantErr: false,
		},
		"valid_no_custom_short": {
			input: &request.ShortenUrlRequest{
				Url:        "http://example.com",
				Expiration: "1h",
			},
			response: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Expiration: time.Hour,
			},
			wantErr: false,
		},
		"valid_no_expiration": {
			input: &request.ShortenUrlRequest{
				Url:         "http://example.com",
				CustomShort: "exmpl",
			},
			response: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour * 24,
			},
			wantErr: false,
		},
		"invalid_exp_time": {
			input: &request.ShortenUrlRequest{
				Url:        "example.com",
				Expiration: "1day",
			},
			wantErr:         true,
			wantErrContains: "failed to parse expiration duration",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := MapShortenUrlRequestToEntity(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErrContains)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.NotEmpty(t, result.Short)
			if tc.input.CustomShort != "" {
				require.Equal(t, tc.response.Short, result.Short)
			}
			require.Equal(t, tc.response.BaseUrl, result.BaseUrl)
			require.Equal(t, tc.response.Expiration, result.Expiration)
		})
	}
}

func TestMapShortenentityToResponse(t *testing.T) {
	testCases := map[string]struct {
		input    *entity.ShortenUrl
		expected *response.ShortenUrlResponse
	}{
		"valid_entity": {
			input: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			expected: &response.ShortenUrlResponse{
				Url:         "http://example.com",
				CustomShort: "exmpl",
				Expiration:  "1h0m0s",
			},
		},
		"valid_entity_no_custom_short": {
			input: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Expiration: time.Hour,
			},
			expected: &response.ShortenUrlResponse{
				Url:        "http://example.com",
				Expiration: "1h0m0s",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := MapShortenentityToResponse(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
