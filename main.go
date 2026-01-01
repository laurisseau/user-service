package main

import (
    "log"
	"net/http"
	"github.com/laurisseau/user-service/handler"
    "github.com/laurisseau/user-service/authenticator"
	"github.com/gin-gonic/gin"
    
)

func main() {


    r := gin.Default()


    // Initialize Authenticator
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize Authenticator: %v", err)
	}

	handler.Router(auth, r)

    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to Sportsify user-service",
        })
    })

    r.Run(":8082") // Starts server on http://localhost:8080 user application port
}

