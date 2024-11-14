package repository

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

var (
	errRedis = fmt.Errorf("error connecting to redis")
)

func TestCheckShortInUse(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewUrlRepository(db)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.Background())
	ctx := &gin.Context{Request: req}

	tests := map[string]struct {
		short     string
		mockSetup func()
		expected  bool
		expectErr bool
	}{
		"url_not_exists": {
			short: "nonexistent",
			mockSetup: func() {
				mock.ExpectGet("nonexistent").RedisNil()
			},
			expected:  false,
			expectErr: false,
		},
		"url_exists": {
			short: "exists",
			mockSetup: func() {
				mock.ExpectGet("exists").SetVal("http://example.com")
			},
			expected:  true,
			expectErr: false,
		},
		"redis_error": {
			short: "error",
			mockSetup: func() {
				mock.ExpectGet("error").SetErr(errRedis)
			},
			expected:  false,
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			result, err := repo.CheckShortInUse(ctx, tc.short)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSaveUrl(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewUrlRepository(db)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.Background())
	ctx := &gin.Context{Request: req}

	tests := map[string]struct {
		short     *entity.ShortenUrl
		mockSetup func()
		expectErr bool
	}{
		"save_success": {
			short: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			mockSetup: func() {
				mock.ExpectSet("exmpl", "http://example.com", time.Hour).SetVal("OK")
			},
			expectErr: false,
		},
		"save_fail": {
			short: &entity.ShortenUrl{
				BaseUrl:    "http://example.com",
				Short:      "exmpl",
				Expiration: time.Hour,
			},
			mockSetup: func() {
				mock.ExpectSet("exmpl", "http://example.com", time.Hour).SetErr(fmt.Errorf("error saving URL"))
			},
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			result, err := repo.SaveUrl(ctx, tc.short)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.short, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUrl(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewUrlRepository(db)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.Background())
	ctx := &gin.Context{Request: req}

	testCases := map[string]struct {
		req       string
		mockSetup func()
		res       string
		expectErr bool
	}{
		"url_exists": {
			req: "exists",
			mockSetup: func() {
				mock.ExpectGet("exists").SetVal("http://example.com")
			},
			res:       "http://example.com",
			expectErr: false,
		},
		"url_not_exists": {
			req: "nonexistent",
			mockSetup: func() {
				mock.ExpectGet("nonexistent").RedisNil()
			},
			expectErr: true,
		},
		"redis_error": {
			req: "error",
			mockSetup: func() {
				mock.ExpectGet("error").SetErr(errRedis)
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()
			result, err := repo.GetUrl(ctx, tc.req)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.res, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
