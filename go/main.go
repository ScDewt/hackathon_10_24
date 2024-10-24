package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"fmt"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		conn, err := pgx.Connect(context.Background(), "postgresql://hackathon:hackathon@postgresql:5432/hackathon")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Database is not accessible", "error": err.Error()})
			return
		}
		defer conn.Close(context.Background())

		c.JSON(http.StatusOK, gin.H{"status": "Database is accessible"})
	})

	r.Run(":8080")
}
