package utils

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"encoding/json"
	"fmt"
)

// RedirectIfNoProfile checks if the profile is nil and redirects to /login if it is.
func RedirectIfNoProfile(ctx *gin.Context, profile interface{}) bool {
	if profile == nil {
		log.Println("No profile in session. Redirecting to /login")
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return true
	}
	return false
}

// GetProfileFromSession retrieves the user profile from the session.
func GetProfileFromSession(ctx *gin.Context) interface{} {
	session := sessions.Default(ctx)
	return session.Get("profile")
}

func GetProfileIdFromSession(ctx *gin.Context) string {
	profile := GetProfileFromSession(ctx)
	if profile == nil {
		return ""
	}

	if profileMap, ok := profile.(map[string]interface{}); ok {
		if id, exists := profileMap["sub"]; exists {
			return id.(string)
		}
	}

	return ""
}

func StringToJSON(str string) (map[string]interface{}, error) {
	var userJSON map[string]interface{}
	if err := json.Unmarshal([]byte(str), &userJSON); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return userJSON, nil
}
