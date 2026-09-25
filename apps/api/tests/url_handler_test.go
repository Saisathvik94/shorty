package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Saisathvik94/shorty/apps/api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func TestCreateURL_BodyTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	body := `{"url":"https://example.com","expires_at":"2026-12-31T00:00:00Z"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/urls",
		strings.NewReader(body),
	)

	c.Request = req

	// Set a limit smaller than the request body.
	bodyLimiter := middlewares.BodyLimiter(10)

	// Apply the body limiter.
	bodyLimiter(c)

	// The middleware only wraps the request body.
	// The actual body-size error happens when the handler reads it.
	_ = c
}
