package profile

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/laurisseau/sportsify-config"
	"github.com/laurisseau/user-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/laurisseau/user-service/models"
	"github.com/laurisseau/user-service/authenticator"
	"github.com/laurisseau/user-service/auth0client"
)

	
// Handler to get profile information from the session.
func Handler(ctx *gin.Context) {

	profile := utils.GetProfileFromSession(ctx)

	if utils.RedirectIfNoProfile(ctx, profile) {
		return
	}

	tokenData, err := authenticator.GetManagementAPIAccessToken()

	if err != nil {
		log.Fatal("Failed to get token:", err)
	}

	userId := utils.GetProfileIdFromSession(ctx)
	
	if userId == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userData, err := auth0client.GetUserByID(userId, config.Secrets["AUTH0_DOMAIN"], tokenData)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	userJSON, err := utils.StringToJSON(userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"profile": userJSON,
	})
}

func UpdateHandler(ctx *gin.Context) {

	config.LoadSecretsEnv()

	profile := utils.GetProfileFromSession(ctx)

	if utils.RedirectIfNoProfile(ctx, profile) {
		return
	}

	var updateReq models.Profile

	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("Invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	tokenData, err := authenticator.GetManagementAPIAccessToken()

	if err != nil {
		log.Fatal("Failed to get token:", err)
	}

	userId := utils.GetProfileIdFromSession(ctx)
	
	if userId == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Construct user update payload
	updatePayload := map[string]interface{}{
		"email":         updateReq.Email,
		"name":          updateReq.Name,
		"user_metadata": updateReq.UserMetadata,
	}

	jsonPayload, _ := json.Marshal(updatePayload)

	userData, err := auth0client.UpdateUserProfile(ctx, updateReq, userId, tokenData, jsonPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	userJSON, err := utils.StringToJSON(userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"profile": userJSON,
	})
}





