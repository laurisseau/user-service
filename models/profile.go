package models

type Profile struct {
	EmailVerified bool                   `json:"email_verified"`
	Name          string                 `json:"name"`
	Email         string                 `json:"email"`
	UserMetadata  map[string]interface{} `json:"user_metadata"`
}
