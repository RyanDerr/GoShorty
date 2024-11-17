package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/internal/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
)

var (
	errRedis = fmt.Errorf("error connecting to redis")
)

func TestShortenUrl(t *testing.T) {
	mock, ctx, service := setupMocks(http.MethodPost)

	testCases := map[string]struct {
		request   *entity.ShortenUrl
		response  *entity.ShortenUrl
		mockSetup func()
		wantErr   bool
		status    int
	}{
		"shorten_success": {
			request: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			response: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			mockSetup: func() {
				mock.ExpectGet("exmpl").RedisNil()
				mock.ExpectSet("exmpl", "http://example.com", time.Hour).SetVal("OK")
			},
			wantErr: false,
			status:  http.StatusCreated,
		},
		"shorten_conflict": {
			request: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			mockSetup: func() {
				mock.ExpectGet("exmpl").SetVal("http://example.com")
			},
			wantErr: true,
			status:  http.StatusConflict,
		},
		"redis_get_error": {
			request: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			mockSetup: func() {
				mock.ExpectGet("exmpl").SetErr(errRedis)
			},
			wantErr: true,
			status:  http.StatusInternalServerError,
		},
		"redis_set_error": {
			request: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			mockSetup: func() {
				mock.ExpectGet("exmpl").RedisNil()
				mock.ExpectSet("exmpl", "http://example.com", time.Hour).SetErr(errRedis)
			},
			wantErr: true,
			status:  http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			resp, status, err := service.ShortenUrl(ctx, tc.request)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Equal(t, tc.status, status)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tc.response, resp)
		})
	}
}

func TestResolveUrl(t *testing.T) {
	mock, ctx, service := setupMocks(http.MethodGet)

	testCases := map[string]struct {
		short     string
		response  string
		mockSetup func()
		wantErr   bool
		status    int
	}{
		"resolve_success": {
			short:    "exmpl",
			response: "http://example.com",
			mockSetup: func() {
				mock.ExpectGet("exmpl").SetVal("http://example.com")
			},
			wantErr: false,
			status:  http.StatusOK,
		},
		"resolve_not_found": {
			short: "exmpl",
			mockSetup: func() {
				mock.ExpectGet("exmpl").RedisNil()
			},
			wantErr: true,
			status:  http.StatusNotFound,
		},
		"redis_get_error": {
			short: "exmpl",
			mockSetup: func() {
				mock.ExpectGet("exmpl").SetErr(errRedis)
			},
			wantErr: true,
			status:  http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			resp, status, err := service.ResolveUrl(ctx, tc.short)
			if tc.wantErr {
				require.Error(t, err)
				require.Empty(t, resp)
				require.Equal(t, tc.status, status)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.response, resp)
		})
	}
}

func setupMocks(httpMethod string) (redismock.ClientMock, *gin.Context, *ShortenUrlService) {
	db, mock := redismock.NewClientMock()
	repo := repository.NewUrlRepository(db)
	service := NewShortenUrlService(repo)
	req, _ := http.NewRequest(httpMethod, "/", nil)
	req = req.WithContext(context.Background())
	ctx := &gin.Context{Request: req}

	return mock, ctx, service
}
