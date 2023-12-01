package routes

import (
	"github.com/gin-gonic/gin"
	Controllers "github.com/oxbase/portfolio_svc/pkg/controllers"
	"gorm.io/gorm"
)

var PortfolioRoutes = func(router *gin.Engine, db *gorm.DB) {
	router.GET("/healthy", Controllers.HealthCheck)

	v1 := router.Group("/api/v1")

	v1.GET("/test", func(c *gin.Context) { Controllers.TestController(c, db) })

	v1.GET("/solana/portfolio/:sol-address", func(c *gin.Context) { Controllers.SolanaController(c, db) })
}
