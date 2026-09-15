package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Shorty API")
	})
	log.Println("API Running on http://localhost:3000")

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
