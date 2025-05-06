package response

import "github.com/gin-gonic/gin"

// Standard HTTP status codes
const (
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusUnprocessableEntity = 422
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusAccepted            = 202
	StatusConflict            = 409
)

// Standard error messages
const (
	MsgInvalidInput       = "Invalid input format"
	MsgEmptyPassword      = "Password cannot be empty"
	MsgUserCreationFailed = "Failed to create user"
)

// Video error messages
const (
	MsgVideoCreationFailed      = "Failed to create video entry"
	MsgVideoUploadFailed        = "Failed to upload video"
	MsgVideoNotFound            = "Video not found"
	MsgVideoAlreadyExists       = "Video already exists"
	MsgVideoUploadNotAllowed    = "Video upload not allowed"
	MsgVideoURLGenerationFailed = "Failed to generate video upload URL"
	MsgDuplicateVideoTitle      = "The video title already exists."
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error sends a standardized error response to the client
// Parameters:
//   - c: The Gin context which manages the HTTP request and response
//   - status: HTTP status code to return (e.g., 400, 404, 500)
//   - message: Short description of the error suitable for end users
//   - details: Additional information about the error, useful for debugging
//
// Returns: Sends a JSON response with the error information and sets the HTTP status code
func Error(c *gin.Context, status int, message string, details string) {
	c.JSON(status, ErrorResponse{
		Status:  status,
		Message: message,
		Details: details,
	})
}

// Success sends a standardized success response
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
