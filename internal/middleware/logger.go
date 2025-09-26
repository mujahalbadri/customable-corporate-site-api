package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin.HanlderFunc that logs requests and responses
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string

		if param.IsOutputColor() {
			statusColor = getStatusColor(param.StatusCode)
			methodColor = getMethodColor(param.Method)
			resetColor = "\033[0m"
		}

		// Custom log format
		return fmt.Sprintf("%s[CORPORATE-API]%s %s - [%s] %s\"%s %s %s %d%s\" %s %s \"%s\" \"%s\"\n",
			"\033[90m", resetColor, // Gray for app name
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			methodColor, param.Method, resetColor,
			param.Path,
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.Request.UserAgent(),
			param.Request.Header.Get("X-Request-ID"),
			param.ErrorMessage,
		)
	})
}

// StructuredLogger returns a structured logger middleware
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		param := gin.LogFormatterParams{
			Request:      c.Request,
			TimeStamp:    time.Now(),
			Latency:      time.Since(start),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     c.Writer.Size(),
		}

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Custom structured log
		logStructuredRequest(param)
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Request.Header.Set("X-Request-ID", requestID)
		}

		// Set the request ID in the response header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingWithConfig returns logger middleware with custom configuration
func LoggingWithConfig(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Skip logging for specified paths
		for _, p := range config.SkipPaths {
			if p == path {
				c.Next()
				return
			}
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Skip if latency is below threshold
		if config.MinLatency > 0 && latency < config.MinLatency {
			return
		}

		param := gin.LogFormatterParams{
			Request:      c.Request,
			TimeStamp:    time.Now(),
			Latency:      latency,
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     c.Writer.Size(),
		}

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		if config.CustomFormatter != nil {
			fmt.Print(config.CustomFormatter(param))
		} else {
			logStructuredRequest(param)
		}
	}
}

// LoggingConfig represents configuration for the logging middleware
type LoggerConfig struct {
	SkipPaths       []string
	MinLatency      time.Duration
	CustomFormatter func(param gin.LogFormatterParams) string
}

// Helper functions for colors and logging
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[97;42m" // White text, green background
	case code >= 300 && code < 400:
		return "\033[90;47m" // Black text, white background
	case code >= 400 && code < 500:
		return "\033[90;43m" // Black text, yellow background
	default:
		return "\033[97;41m" // White text, red background
	}
}
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[94m" // Blue
	case "POST":
		return "\033[92m" // Green
	case "PUT":
		return "\033[93m" // Yellow
	case "DELETE":
		return "\033[91m" // Red
	case "PATCH":
		return "\033[95m" // Magenta
	case "HEAD":
		return "\033[96m" // Cyan
	case "OPTIONS":
		return "\033[90m" // Gray
	default:
		return "\033[97m" // White
	}
}

func logStructuredRequest(param gin.LogFormatterParams) {
	statusEmoji := getStatusEmoji(param.StatusCode)
	methodEmoji := getMethodEmoji(param.Method)

	fmt.Printf("%s %s %s | %s | %v | %s | %d bytes | %s\n",
		statusEmoji,
		methodEmoji,
		param.Method,
		param.Path,
		param.Latency,
		param.ClientIP,
		param.BodySize,
		param.TimeStamp.Format("15:04:05"),
	)

	// Log error if present
	if param.ErrorMessage != "" {
		fmt.Printf("   ❌ Error: %s\n", param.ErrorMessage)
	}
}

func getStatusEmoji(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "✅"
	case code >= 300 && code < 400:
		return "🔄"
	case code >= 400 && code < 500:
		return "⚠️"
	default:
		return "❌"
	}
}

func getMethodEmoji(method string) string {
	switch method {
	case "GET":
		return "📖"
	case "POST":
		return "📝"
	case "PUT":
		return "✏️"
	case "DELETE":
		return "🗑️"
	case "PATCH":
		return "🔧"
	case "HEAD":
		return "👁️"
	case "OPTIONS":
		return "❓"
	default:
		return "📡"
	}
}

func generateRequestID() string {
	// Simple unique ID generator (for demonstration purposes)
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
