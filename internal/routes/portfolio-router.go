package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/0xbase-Corp/portfolio_svc/internal/controllers"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/bitcoin"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/debank"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/solana"
)

var PortfolioRoutes = func(router *gin.Engine, db *gorm.DB) {
	router.GET("/healthy", controllers.HealthCheck)

	bitcoinAPIClient := &bitcoin.BitcoinAPI{}
	solanaAPIClient := &solana.SolanaAPI{}
	debankAPIClient := &debank.DebankAPI{}

	v1 := router.Group("/api/v1")

	v1.GET("/test", func(c *gin.Context) { controllers.TestController(c, db) })

	v1.GET("/portfolio/solana", func(c *gin.Context) { controllers.SolanaController(c, db, solanaAPIClient) })

	v1.GET("/portfolio/solana-wallet/:wallet-id", func(c *gin.Context) { controllers.GetSolanaController(c, db) })

	v1.GET("/portfolio/btc", func(c *gin.Context) { controllers.BitcoinController(c, db, bitcoinAPIClient) })

	v1.GET("/portfolio/btc-wallet/:wallet-id", func(c *gin.Context) { controllers.GetBtcDataController(c, db) })

	v1.GET("/portfolio/debank", func(c *gin.Context) { controllers.DebankController(c, db, debankAPIClient) })

	v1.POST("/all-portfolio", func(c *gin.Context) { controllers.AllPortfolioController(c, db) })

}
