package main

import (
	"log"

	docs "github.com/0xbase-Corp/portfolio_svc/cmd/docs"
	"github.com/0xbase-Corp/portfolio_svc/pkg/configs"
	"github.com/0xbase-Corp/portfolio_svc/pkg/routes"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r := gin.Default()
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
