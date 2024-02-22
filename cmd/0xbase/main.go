package main

import (
	"log"
	"sync"

	"github.com/0xbase-Corp/portfolio_svc/docs"
	"github.com/0xbase-Corp/portfolio_svc/internal/middlewares"
	"github.com/0xbase-Corp/portfolio_svc/internal/routes"
	"github.com/0xbase-Corp/portfolio_svc/shared/configs"
	"github.com/0xbase-Corp/portfolio_svc/shared/migrations"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	once sync.Once
)

//	@title			0xBase-Corp API
//	@version		1.0
//	@description	This is Portfolio server API documentation.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	0xSantos
//	@contact.url	http://www.0xbase.org
//	@contact.email	help@0xbase.org

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5050
//	@BasePath	/api/v1

func main() {
	//Loading Environment variables from app.env
	configs.InitEnvConfigs()

	db := configs.GetDB()

	// This ensures that the migration process is executed only once, regardless of how many times main() is called.
	// Execute database migration once
	once.Do(func() {
		log.Println("Starting database migration...")
		err := migrations.Migrate(db) // Call the Migrate function from migrations package
		if err != nil {
			log.Fatalf("Database migration failed: %v", err)
		}
		log.Println("Database migration completed successfully.")
	})

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	docs.SwaggerInfo.BasePath = "/api/v1"

	//gin warning: "you trusted all proxies this is not safe. we recommend you to set a value"
	r.ForwardedByClientIP = true
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Failed to setup trusted Proxies")
	}

	routes.PortfolioRoutes(r, db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := r.Run(configs.EnvConfigVars.Port); err != nil {
		log.Println("Server failed to start ", err)
	}
}
