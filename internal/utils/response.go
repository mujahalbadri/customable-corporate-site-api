package utils

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// API Response represents a standard API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// Pagination Response represents paginated data
type PaginationResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

// Error Detail represents a single error detail
type ErrorDetail struct {
	Code    string      `json:"code,omitempty"`
	Field   string      `json:"field,omitempty"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// Validation Error represents a validation error structure
type ValidationError struct {
	Message string        `json:"message"`
	Error   []ErrorDetail `json:"errors"`
}

// Success Response sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	c.JSON(statusCode, response)
}

// Error Response sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	var errorData string
	if err != nil {
		errorData = err.Error()
	}

	response := APIResponse{
		Success:   false,
		Message:   message,
		Error:     errorData,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	c.JSON(statusCode, response)
}

// Validation Error Response sends a validation error response
func ValidationErrorResponse(c *gin.Context, statusCode int, message string, errors []ErrorDetail) {
	ValidationError := ValidationError{
		Message: message,
		Error:   errors,
	}

	response := APIResponse{
		Success:   false,
		Message:   "Validation Error",
		Data:      ValidationError,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	c.JSON(statusCode, response)
}

// Paginated Success Response sends a paginated success response
func PaginatedSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination Pagination) {
	response := PaginationResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
		RequestID:  getRequestID(c),
	}

	c.JSON(statusCode, response)
}

// NotFoundResponse sends a 404 not found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, 404, message+" not found", nil)
}

// UnauthorizedResponse sends a 401 unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, 401, message, nil)
}

// ForbiddenResponse sends a 403 forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponse(c, 403, message, nil)
}

// BadRequestResponse sends a 400 bad request response
func BadRequestResponse(c *gin.Context, message string, err error) {
	if message == "" {
		message = "Bad Request"
	}
	ErrorResponse(c, 400, message, err)
}

// InternalServerErrorResponse sends a 500 internal server error response
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	if message == "" {
		message = "Internal Server Error"
	}
	ErrorResponse(c, 500, message, err)
}

// ConflictResponse sends a 409 conflict response
func ConflictResponse(c *gin.Context, message string, err error) {
	if message == "" {
		message = "Conflict"
	}
	ErrorResponse(c, 409, message, err)
}

// CreatedResponse sends a 201 created response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Resource created successfully"
	}
	SuccessResponse(c, 201, message, data)
}

// NoContentResponse sends a 204 no content response
func NoContentResponse(c *gin.Context) {
	c.Status(204)
}

// HealthCheckResponse sends a health check response
func HealthCheckResponse(c *gin.Context, status string, details interface{}) {
	response := gin.H{
		"status":    status,
		"timestamp": time.Now(),
		"service":   getEnv("APP_NAME", "Customable Corporate Site API"),
		"version":   "1.0.0",
		"details":   details,
	}

	if details != nil {
		response["details"] = details
	}

	statusCode := 200
	if status != "healthy" && status != "ok" {
		statusCode = 503
	}

	c.JSON(statusCode, response)
}

// Helper function to get request ID from context
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if idStr, ok := requestID.(string); ok {
			return idStr
		}
	}
	return ""
}

// Calculate Pagination metadata
func CalculatePagination(currentPage, pageSize int, totalItems int) Pagination {
	totalPages := (totalItems + pageSize - 1) / pageSize // Ceiling division
	return Pagination{
		TotalItems:  totalItems,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}
}

// Parse Pagination extracts pagination parameters from the request
func ParsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if pageNum := parseInt(p); pageNum > 0 {
			page = pageNum
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if pageSizeNum := parseInt(ps); pageSizeNum > 0 && pageSizeNum <= 100 {
			pageSize = pageSizeNum
		}
	}

	return page, pageSize
}

// parseInt is a helper function to parse string to int with error handling
func parseInt(s string) int {
	if s == "" {
		return 0
	}

	// simple conversion, can be enhanced with more error handling
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}

// Response With Metadata sends a response with additional metadata
func ResponseWithMetadata(c *gin.Context, statusCode int, message string, data interface{}, metadata interface{}) {
	response := map[string]interface{}{
		"success":    statusCode >= 200 && statusCode < 300,
		"message":    message,
		"data":       data,
		"metadata":   metadata,
		"timestamp":  time.Now(),
		"request_id": getRequestID(c),
	}
	c.JSON(statusCode, response)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
