package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/Saisathvik94/shorty/apps/api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	middleware := middlewares.SecurityHeaders()

	c.Request = httptest.NewRequest("GET", "/health", nil)

	middleware(c)

	if got := c.Writer.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("expected X-Frame-Options to be DENY, got %q", got)
	}

	if got := c.Writer.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("expected X-Content-Type-Options to be nosniff, got %q", got)
	}
}
