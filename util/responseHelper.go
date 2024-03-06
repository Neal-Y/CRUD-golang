package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// status code just an integer
func ResponseWithJson(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, gin.H{"payload": payload})
}

// error handler
func ResponseWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		ResponseWithError(c, http.StatusBadRequest, "Failed to read body")
		return false
	}

	if err := validate.Struct(obj); err != nil {
		ResponseWithError(c, http.StatusBadRequest, err.Error())
		return false
	}

	return true
}
