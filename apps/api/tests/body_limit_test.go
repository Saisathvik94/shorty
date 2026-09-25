package tests

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Saisathvik94/shorty/apps/api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func TestBodyLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        string
		maxBodySize int64
		wantErr     bool
	}{
		{
			name:        "body within limit",
			body:        strings.Repeat("a", 100),
			maxBodySize: 200,
			wantErr:     false,
		},
		{
			name:        "body exceeds limit",
			body:        strings.Repeat("a", 200),
			maxBodySize: 100,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(
				"POST",
				"/api/urls",
				strings.NewReader(tt.body),
			)

			c.Request = req

			middleware := middlewares.BodyLimiter(tt.maxBodySize)

			middleware(c)

			data, err := io.ReadAll(c.Request.Body)

			if tt.wantErr && err == nil {
				t.Fatal("expected body size error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !tt.wantErr && string(data) != tt.body {
				t.Fatalf("expected body %q, got %q", tt.body, string(data))
			}
		})
	}
}
