package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupCORSTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/urls", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "created",
		})
	})

	return r
}

func TestCORS_AllowedOrigin(t *testing.T) {
	router := setupCORSTestRouter()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/urls",
		nil,
	)

	req.Header.Set("Origin", "http://localhost:5173")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	got := recorder.Header().Get("Access-Control-Allow-Origin")

	if got != "http://localhost:5173" {
		t.Fatalf(
			"expected Access-Control-Allow-Origin %q, got %q",
			"http://localhost:5173",
			got,
		)
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	router := setupCORSTestRouter()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/urls",
		nil,
	)

	req.Header.Set("Origin", "http://localhost:3000")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	got := recorder.Header().Get("Access-Control-Allow-Origin")

	if got != "" {
		t.Fatalf(
			"expected no Access-Control-Allow-Origin header, got %q",
			got,
		)
	}
}

func TestCORS_Preflight(t *testing.T) {
	router := setupCORSTestRouter()

	req := httptest.NewRequest(
		http.MethodOptions,
		"/api/urls",
		nil,
	)

	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			recorder.Code,
		)
	}

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf(
			"expected allowed origin %q, got %q",
			"http://localhost:5173",
			got,
		)
	}

	if got := recorder.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("expected Access-Control-Allow-Methods header")
	}
}
