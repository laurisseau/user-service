package main

import (
    "log"
	"net/http"
	"sportsify/user-service/handler"
    "sportsify/user-service/authenticator"
	"github.com/gin-gonic/gin"
    "github.com/laurisseau/sportsify-config"
)

func main() {

    r := gin.Default()

    db := config.DB()

    // Initialize Authenticator
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize Authenticator: %v", err)
	}

	handler.Router(db, auth, r)

    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to Sportsify!",
        })
    })

    r.Run(":8080") // Starts server on http://localhost:8080 user application port
}

