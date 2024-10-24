package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"healthcheck/source_manager"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	// Настройки CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // Разрешенные домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Разрешенные методы
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	dsn := getEnv("DSN", "postgresql://hackathon:hackathon@postgresql:5432/hackathon")

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	r.GET("/api/health", func(c *gin.Context) {
		_, err := conn.Exec(c.Request.Context(), "SELECT 1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Database is not accessible", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Database is accessible"})
	})

	source_manager.NewSourceManagerHandler(conn).RegisterHandler(r)

	r.Run(":8080")
}

// getEnv получает строку из переменных окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
