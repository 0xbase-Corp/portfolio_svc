package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/0xbase-Corp/portfolio_svc/internal/controllers"
)

var PortfolioRoutes = func(router *gin.Engine, db *gorm.DB) {
	router.GET("/healthy", controllers.HealthCheck)

	v1 := router.Group("/api/v1")

	v1.GET("/test", func(c *gin.Context) { controllers.TestController(c, db) })

	v1.GET("/portfolio/solana/:sol-address", func(c *gin.Context) { controllers.SolanaController(c, db) })

	v1.GET("/portfolio/btc/:btc-address", func(c *gin.Context) { controllers.BitcoinController(c, db) })

	v1.GET("/portfolio/debank/:debank-address", func(c *gin.Context) { controllers.DebankController(c, db) })
}
