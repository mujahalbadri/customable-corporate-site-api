package services

import (
	"customable-corporate-site-api/internal/models"
	"customable-corporate-site-api/internal/repositories/interfaces"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthService defines the interface for authentication services.
type AuthService struct {
	userRepo  interfaces.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// JWT Claims structure
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Request DTOs
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
}

// Response DTOs
type TokenResponse struct {
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	ExpiresIn    int64                `json:"expires_in"`
	TokenType    string               `json:"token_type"`
	User         *models.UserResponse `json:"user"`
}

type AuthResponse struct {
	Message string               `json:"message"`
	User    *models.UserResponse `json:"user"`
	Token   *TokenResponse       `json:"token,omitempty"`
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userRepo interfaces.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	newUser := &models.User{
		Email:     req.Email,
		Password:  req.Password, // Assume password is hashed in repository layer
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
		IsActive:  true,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, errors.New("failed to create user account")
	}

	return &AuthResponse{
		Message: "User registered successfully",
		User:    newUser.ToResponse(),
	}, nil
}

// Login authenticates a user and returns JWT tokens.
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Fetch user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT tokens
	tokenResponse, err := s.generateTokenResponse(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &AuthResponse{
		Message: "Login successful",
		User:    user.ToResponse(),
		Token:   tokenResponse,
	}, nil
}

// RefreshToken generates a new access token using a refresh token.
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// Parse and validate the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Verify is a refresh token
	if claims.Subject != "refresh_token" {
		return nil, errors.New("invalid token type")
	}

	// Get user to ensure they still exist and are active
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	return s.generateTokenResponse(user)
}

// GetProfile retrieves the profile of the authenticated user.
func (s *AuthService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

// UpdateProfile updates the profile of the authenticated user.
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = strings.TrimSpace(req.FirstName)
	}

	if req.LastName != "" {
		user.LastName = strings.TrimSpace(req.LastName)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user.ToResponse(), nil
}

// ValidateToken validates a JWT token and returns the associated user.
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Fetch user by ID
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	return claims, nil
}

// Private helper methods

// generateTokenResponse creates access and refresh tokens for a user.
func (s *AuthService) generateTokenResponse(user *models.User) (*TokenResponse, error) {
	// Create access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
		TokenType:    "Bearer",
		User:         user.ToResponse(),
	}, nil
}

// generateAccessToken creates a JWT access token for a user.
func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
			Issuer:    "customable-corporate-site-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken creates a JWT refresh token for a user.
func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh token valid for 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
			Issuer:    "customable-corporate-site-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
