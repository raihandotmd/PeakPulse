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

// Signup godoc
// @Summary      Register a new user
// @Description  Creates a new user account using Supabase Auth and stores no local password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      models.SignupRequest   true  "Email and password"
// @Success      201      {object}  models.SignupResponse   "Signup successful"
// @Failure      400      {object}  models.ErrorResponse    "Invalid request or duplicate email"
// @Failure      409      {object}  models.ErrorResponse    "Email already exists"
// @Failure      500      {object}  models.ErrorResponse    "Internal server error"
// @Router       /signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var body models.SignupRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	DB := db.GetGormClient()
	// Check for existing email
	var existUser models.User
	if err := DB.Where("email = ?", body.Email).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: "Email already exists"})
		return
	}

	// Attempt Supabase signup
	authResp, err := h.Supa.Auth.Signup(types.SignupRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to sign up user"})
		return
	}

	c.JSON(http.StatusCreated, models.SignupResponse{
		Message: "Sign-up successful. Welcome!",
		User: models.User{
			ID:    authResp.User.ID,
			Email: authResp.User.Email,
		},
	})
}

// Login godoc
// @Summary      Log in and receive a JWT
// @Description  Authenticates a user via Supabase Auth and returns a JWT for API access.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      models.LoginRequest    true  "Email and password"
// @Success      200      {object}  models.LoginResponse   "Login successful"
// @Failure      400      {object}  models.ErrorResponse   "Invalid request payload"
// @Failure      401      {object}  models.ErrorResponse   "Authentication failed"
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var body models.LoginRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	sess, err := h.Supa.Auth.SignInWithEmailPassword(body.Email, body.Password)
	if err != nil {
		// Try to extract Supabase error message
		var supaErr struct {
			Msg string `json:"msg"`
		}
		if idx := strings.Index(err.Error(), "{"); idx != -1 {
			if e := json.Unmarshal([]byte(err.Error()[idx:]), &supaErr); e == nil && supaErr.Msg != "" {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: supaErr.Msg})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: sess.Session.AccessToken})
}
