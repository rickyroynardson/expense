package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rickyroynardson/expense/internal/auth"
	"github.com/rickyroynardson/expense/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatalln("APP_PORT is not set")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatalln("DATABASE_URL is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalln("JWT_SECRET is not set")
	}
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		log.Fatalln("COOKIE_DOMAIN is not set")
	}

	cfg := &utils.Config{
		JwtSecret:    jwtSecret,
		CookieDomain: cookieDomain,
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer conn.Close(context.Background())

	authRepository := auth.NewRepository(conn)
	authHandler := auth.NewHandler(cfg, authRepository)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	auth := router.Group("/api/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)

	if err := router.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
