package middlewares

import (
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware for handling errors and responding with appropriate JSON.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			handleError(c, err)
		}
	}
}

// handleError handles a specific error, sending the appropriate JSON response.
func handleError(c *gin.Context, err *gin.Error) {
	switch e := err.Err.(type) {
	case *errors.APIError:
		c.AbortWithStatusJSON(e.Code, e)
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Service Unavailable"})
	}
}
