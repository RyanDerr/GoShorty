package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckURLInUse(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewURLRepository(db)
	ctx := &gin.Context{}

	tests := map[string]struct {
		id        string
		dbResult  string
		expectErr bool
	}{
		"url_in_use":     {generateUUID(), "http://example.com", true},
		"url_not_in_use": {generateUUID(), "", false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectGet(tc.id).SetVal(tc.dbResult)
			err := repo.CheckURLInUse(ctx, tc.id)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, "url already in use", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveURL(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewURLRepository(db)
	ctx := &gin.Context{}

	tests := map[string]struct {
		id         string
		url        string
		expiration time.Duration
		expectErr  bool
	}{
		"happy_path_save_url": {generateUUID(), "http://example.com", 24 * time.Hour, false},
		"fail_to_save_url":    {generateUUID(), "http://example.com", 24 * time.Hour, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectErr {
				mock.ExpectSet(tc.id, tc.url, tc.expiration).SetErr(fmt.Errorf("failed to save URL"))
			} else {
				mock.ExpectSet(tc.id, tc.url, tc.expiration).SetVal("OK")
			}
			err := repo.SaveURL(ctx, tc.id, tc.url, tc.expiration)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, "failed to save URL", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewURLRepository(db)
	ctx := &gin.Context{}

	tests := map[string]struct {
		id        string
		dbResult  string
		expectErr bool
	}{
		"happy_path_retrieve_url": {generateUUID(), "http://example.com", false},
		"url_not_found":           {generateUUID(), "", true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.dbResult == "" {
				mock.ExpectGet(tc.id).RedisNil()
			} else {
				mock.ExpectGet(tc.id).SetVal(tc.dbResult)
			}
			url, err := repo.GetURL(ctx, tc.id)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, "", url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.dbResult, url)
			}
		})
	}
}

func generateUUID() string {
	uuid := uuid.New()
	return uuid.String()[:6]
}
