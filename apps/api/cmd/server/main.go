package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Saisathvik94/shorty/apps/api/internal/database"
	"github.com/Saisathvik94/shorty/apps/api/internal/handlers"
	"github.com/Saisathvik94/shorty/apps/api/internal/repository"
	"github.com/Saisathvik94/shorty/apps/api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	// health route
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Shorty API")
	})

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	// Connection pool
	pool, err := database.PostgresPool(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// create repository
	repo := repository.NewURLRepository(pool)

	// create service
	service := services.NewURLService(repo)

	// create handler
	handler := handlers.NewURLHandler(service)

	// Routes
	r.POST("/api/urls", handler.CreateURL)
	r.GET("/:shortCode", handler.Redirect)

	// Graceful Shutdown
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}
	// A goroutine with no name
	go func() {
		log.Printf("API Running on: http://localhost%s\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // SIGINT (Ctrl + C) and SIGTERM (Docker*/ Kubernetes)

	<-quit // wait

	log.Println("Shutting Down the server....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Forced to ShutDown: %v", err)
	}

	log.Println("Server Stopped")
}
