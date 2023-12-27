package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@BasePath	/api/v1

// PingExample godoc
//
//	@Summary		ping example
//	@Schemes		models.BtcChainAPI
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string} Helloworld
//	@Router			/portfolio/btc/:btc-address [get]
func BitcoinController(c *gin.Context, db *gorm.DB) {
	btcAddress := c.Param("btc-address")

	url := "https://chain.api.btc.com/v3/address/" + btcAddress

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	//send the request
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("body ", body)
	var btcResponse models.BtcChainAPI
	if err := json.Unmarshal(body, &btcResponse); err != nil {
		log.Fatal(err)
	}

	fmt.Println("BTC PResponse ", btcResponse)
	//TODO: add the data into database

	c.JSON(http.StatusOK, gin.H{
		"data": btcResponse,
	})
}
