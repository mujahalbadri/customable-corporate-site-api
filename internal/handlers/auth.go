package handlers

import (
	"customable-corporate-site-api/internal/services"
	"customable-corporate-site-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration.
// @Summary Register a new user
// @Description Register a new user with email, password, first name, and last name.
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body services.RegisterRequest true "Register Request"
// @Success 201 {object} services.AuthResponse
// @Failure 400 {object} services.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Call service to register user
	resp, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", resp)
}

// Login hanldles user login.
// @Summary Login a user
// @Description Login a user with email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body services.LoginRequest true "Login Request"
// @Success 200 {object} services.AuthResponse
// @Failure 400 {object} services.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Call service to login user
	resp, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User logged in successfully", resp)
}

// RefreshToken handles token refresh.
// @Summary Refresh JWT tokens
// @Description Refresh JWT access and refresh tokens using a valid refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshTokenRequest body services.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} services.TokenResponse
// @Failure 400 {object} services.ErrorResponse
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Call service to refresh token
	tokenResp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", tokenResp)
}

// GetProfile handles fetching the authenticated user's profile.
// @Summary Get user profile
// @Description Retrieve the profile of the authenticated user.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} services.UserResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 404 {object} services.ErrorResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get UserID from JWT middleware
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Type assertion to uint
	id, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	// Call service to get user profile
	profile, err := h.authService.GetProfile(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", profile)
}

// UpdateProfile handles updating the authenticated user's profile.
// @Summary Update user profile
// @Description Update the profile of the authenticated user.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updateProfileRequest body services.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} services.UserResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 404 {object} services.ErrorResponse
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get UserID from JWT middleware
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Type assertion to uint
	id, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	var req services.UpdateProfileRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Call service to update user profile
	updatedProfile, err := h.authService.UpdateProfile(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", updatedProfile)
}

// GetCurrentUser is an alias for GetProfile to maintain backward compatibility.
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	h.GetProfile(c)
}

// RefreshTokenRequest represents the request payload for refreshing tokens.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
