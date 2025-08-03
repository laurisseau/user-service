package auth0client

import (
	"bytes"
	"io"
	"log"
	"fmt"
	"net/http"
	"github.com/laurisseau/sportsify-config"
	"github.com/gin-gonic/gin"
	"github.com/laurisseau/user-service/models"
)

func UpdateUserProfile(ctx *gin.Context, updateReq models.Profile, userId string, tokenData string, jsonPayload []byte) (string, error){
// Send PATCH request to Auth0
	url := "https://" + config.Secrets["AUTH0_DOMAIN"] + "/api/v2/users/" + userId
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		log.Printf("Error creating PATCH request: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenData)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return string(body), nil
}

// GetUserByID fetches a user from Auth0 by ID and returns the raw JSON string
func GetUserByID(userID, authDomain, token string) (string, error) {
	url := fmt.Sprintf("https://%s/api/v2/users/%s", authDomain, userID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to Auth0 failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth0 returned error response: %s", body)
		return "", fmt.Errorf("auth0 returned status: %d", resp.StatusCode)
	}

	return string(body), nil
}