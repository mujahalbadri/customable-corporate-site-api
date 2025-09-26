package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware to handle Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", getAllowedOrigin(origin))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-CSRF-Token, X-Requested-With, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORS With Config creates a CORS middleware with custom configuration
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		allowedOrigin := "*"
		if len(config.AllowedOrigins) > 0 {
			allowedOrigin = getAllowedOriginFromConfig(origin, config.AllowedOrigins)
		}

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowedMethods, "GET, POST, PUT, DELETE, OPTIONS"))
		c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowedHeaders, "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With, Authorization"))
		c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposedHeaders, "Content-Length, X-CSRF-Token, X-Requested-With, Authorization"))
		c.Header("Access-Control-Allow-Credentials", boolToString(config.AllowCredentials))

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSConfig holds the configuration for the CORS middleware
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// Default CORS configuration returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "X-Requested-With", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "X-CSRF-Token", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// Development CORS configuration returns CORS configuration for development environment
func DevelopmentCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:8080",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "X-Requested-With", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "X-CSRF-Token", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// Helper functions

func getAllowedOrigin(requestOrigin string) string {
	// In development, be more permissive
	developmentOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:8080",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"http://127.0.0.1:8080",
	}

	for _, origin := range developmentOrigins {
		if origin == requestOrigin {
			return origin
		}
	}

	// For production, you should specify allowed origins
	// Here we just return "*" for simplicity
	return "*"
}

func getAllowedOriginFromConfig(origin string, allowedOrigins []string) string {
	for _, allallowedOrigins := range allowedOrigins {
		if allallowedOrigins == "*" {
			return "*"
		}
		if allallowedOrigins == origin {
			return origin
		}
	}
	return "null"
}

func joinStrings(slice []string, defaultValue string) string {
	if len(slice) == 0 {
		return defaultValue
	}

	result := ""
	for i, str := range slice {
		if i > 0 {
			result += ", "
		}
		result += str
	}
	return result
}

func boolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
