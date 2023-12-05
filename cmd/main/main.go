package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oxbase/portfolio_svc/pkg/configs"
	"github.com/oxbase/portfolio_svc/pkg/routes"
)

func main() {

	//Loading Environment variables from app.env
	configs.InitEnvConfigs()

	db := configs.GetDB()

	r := gin.Default()
	//gin warning: "you trusted all proxies this is not safe. we recommend you to set a value"
	r.ForwardedByClientIP = true
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Failed to setup trusted Proxies")
	}

	routes.PortfolioRoutes(r, db)

	if err := r.Run(configs.GetPort()); err != nil {
		log.Println("Server failed to start ", err)
	}
}
