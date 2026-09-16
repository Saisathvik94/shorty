package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Shorty API")
	})

	// Graceful Shutdown
	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	// A goroutine with no name
	go func() {
		log.Println("API Running on: http://localhost:3000")

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
