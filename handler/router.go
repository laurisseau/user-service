package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/laurisseau/user-service/authenticator"
    "encoding/gob"
    "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
    "github.com/laurisseau/user-service/handler/login"
    "github.com/laurisseau/user-service/handler/callback"
    "github.com/laurisseau/user-service/handler/profile"
    "github.com/laurisseau/user-service/handler/middleware"
    "github.com/laurisseau/user-service/handler/logout"
)

// New registers the routes and returns the router.
func Router(auth *authenticator.Authenticator, router *gin.Engine) {
	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

    // Create a /users route group
    userGroup := router.Group("/users")
    
    {
        userGroup.GET("/login", login.Handler(auth))
        userGroup.GET("/callback", callback.Handler(auth))
        userGroup.GET("/logout", logout.Handler)
        userGroup.GET("/profile", middleware.IsAuthenticated, profile.Handler)
        userGroup.PATCH("/profile/update", middleware.IsAuthenticated, profile.UpdateHandler)
    }
}
