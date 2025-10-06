package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents the structure of JWT claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth middleware validates JWT tokens and extracts user information
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required", "message": "Unauthorized"})
			c.Abort()
			return
		}

		// Check bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer <token>'", "message": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": "Token validation failed: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": "Please login again to obtain a new token"})
			c.Abort()
			return
		}

		// Extract claims and set user information in context
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims", "message": "Token format is not recognized"})
			c.Abort()
			return
		}

		// Verify this is an access token
		if claims.Subject != "access_token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type", "message": "Only access tokens are allowed"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RequireRoles middleware checks if the user has one of the required roles
func RequireRoles(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found in token", "message": "Authentication required"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role format", "message": "Authentication required"})
			c.Abort()
			return
		}

		// Check role hierarchy (Admin > Editor > User)
		if !hasRequiredRole(role, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions", "message": "You do not have access to this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware ensures the user has admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRoles("admin")
}

// RequireEditor middleware ensures the user has editor or higher role
func RequireEditor() gin.HandlerFunc {
	return RequireRoles("editor")
}

// hasRequiredRole checks if the user's role meets the required role
func hasRequiredRole(userRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"user":   1,
		"editor": 2,
		"admin":  3,
	}

	userLevel := roleHierarchy[strings.ToLower(userRole)]
	requiredLevel := roleHierarchy[strings.ToLower(requiredRole)]

	return userLevel >= requiredLevel
}

// OptionalAuth middleware allows optional authentication
func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, proceed without authentication
			c.Next()
			return
		}

		// Check Bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || claims.Subject != "access_token" {
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)
		c.Set("authenticated", true)

		c.Next()
	}
}
