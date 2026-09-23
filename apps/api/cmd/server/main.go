package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Saisathvik94/shorty/apps/api/internal/cache"
	"github.com/Saisathvik94/shorty/apps/api/internal/database"
	"github.com/Saisathvik94/shorty/apps/api/internal/handlers"
	"github.com/Saisathvik94/shorty/apps/api/internal/repository"
	"github.com/Saisathvik94/shorty/apps/api/internal/services"
	"github.com/Saisathvik94/shorty/apps/api/internal/workers"
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
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is not set")
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

	// Redis client
	rdb, err := cache.NewRedisClient(redisURL)
	if err != nil {
		log.Fatal(err)
	}

	pingCtx := context.Background()
	pong, err := rdb.Ping(pingCtx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Printf("Successfully connected! Redis responded with: %s\n", pong)

	// passing pool to repository
	repo := repository.NewURLRepository(pool)

	// passing rb to the cache
	cache := cache.NewRedisCache(rdb)

	// passing repo to the services
	service := services.NewURLService(cache, repo)

	// passing the service to the handler
	handler := handlers.NewURLHandler(service)

	// Worker
	w := workers.NewCleanUpWorker(repo, 30*time.Second)
	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()
	go w.Start(appCtx)

	// Routes
	r.POST("/api/urls", handler.CreateURL)
	r.GET("/:shortCode", handler.Redirect)
	r.PUT("/api/urls/:shortCode/deactivate", handler.DeactivateURL)
	r.DELETE("/api/urls/:shortCode", handler.DeleteURL)
	r.PATCH("/api/urls/:shortCode/expiration", handler.UpdateExpiration)

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

	cancelApp()

	defer cancel()
	defer rdb.Close()
	defer pool.Close()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Forced to ShutDown: %v", err)
	}

	log.Println("Server Stopped")
}
