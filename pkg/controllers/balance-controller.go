package controllers

import (
	// "encoding/json"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /db/test/healthy [get]
func TestController(c *gin.Context, db *gorm.DB) {
	users := models.GetAllUsers(db)
	// res, _ := json.Marshal(users)

	c.JSON(http.StatusOK, gin.H{
		"message": users,
	})
}
