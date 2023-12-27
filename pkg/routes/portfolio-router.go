package routes

import (
	Controllers "github.com/0xbase-Corp/portfolio_svc/pkg/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var PortfolioRoutes = func(router *gin.Engine, db *gorm.DB) {
	router.GET("/healthy", Controllers.HealthCheck)

	v1 := router.Group("/api/v1")

	v1.GET("/test", func(c *gin.Context) { Controllers.TestController(c, db) })

	v1.GET("/portfolio/solana/:sol-address", func(c *gin.Context) { Controllers.SolanaController(c, db) })

	v1.GET("/portfolio/btc/:btc-address", func(c *gin.Context) { Controllers.BitcoinController(c, db) })
}
