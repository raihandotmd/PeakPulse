package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raihandotmd/peakPulse/internal/db"
	"github.com/raihandotmd/peakPulse/internal/models"
	"github.com/supabase-community/gotrue-go/types"
	supa "github.com/supabase-community/supabase-go"
)

// AuthHandler holds Supabase client for auth operations
type AuthHandler struct {
	Supa *supa.Client
}

// NewAuthHandler returns a new AuthHandler
func NewAuthHandler(supa *supa.Client) *AuthHandler {
	return &AuthHandler{Supa: supa}
}

// Signup registers a new user
func (h *AuthHandler) Signup(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	DB := db.GetGormDB() // Get the Gorm DB instance
	// Check if the email already exists in the database
	var existUser models.User
	if err := DB.Table("users").Where("email = ?", body.Email).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Attempt to sign up the user
	user, err := h.Supa.Auth.Signup(types.SignupRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up user"})
		return
	}

	// Successful sign-up
	c.JSON(http.StatusCreated, gin.H{
		"message": "Sign-up successful. Welcome!",
		"user":    user.User,
	})
}

// Login authenticates and returns a JWT
func (h *AuthHandler) Login(c *gin.Context) {
	type req struct{ Email, Password string }
	var body req
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sess, err := h.Supa.Auth.SignInWithEmailPassword(body.Email, body.Password)

	if err != nil {
		// Define a struct to parse the error message
		var supabaseError struct {
			Code      int    `json:"code"`
			ErrorCode string `json:"error_code"`
			Msg       string `json:"msg"`
		}

		// Find the JSON part of the error message
		errorMessage := err.Error()
		startIndex := strings.Index(errorMessage, "{")
		if startIndex != -1 {
			jsonPart := errorMessage[startIndex:]

			// Attempt to parse the JSON part
			if jsonErr := json.Unmarshal([]byte(jsonPart), &supabaseError); jsonErr == nil {
				// If parsing is successful, return the error message
				c.JSON(http.StatusUnauthorized, gin.H{"error": supabaseError.Msg})
				return
			}
		}

		// If parsing fails, return the raw error message
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// generate JWT for your API (or return Supabase access token)
	c.JSON(http.StatusOK, gin.H{"token": sess.Session.AccessToken})
}
