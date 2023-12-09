package main

import (
	"log"

	"github.com/gin-gonic/gin"
	docs "github.com/oxbase/portfolio_svc/cmd/docs"
	"github.com/oxbase/portfolio_svc/pkg/configs"
	"github.com/oxbase/portfolio_svc/pkg/controllers"
	"github.com/oxbase/portfolio_svc/pkg/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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

	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", controllers.Helloworld)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes.PortfolioRoutes(r, db)

	if err := r.Run(configs.EnvConfigVars.Port); err != nil {
		log.Println("Server failed to start ", err)
	}
}
