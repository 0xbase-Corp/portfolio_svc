package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes models.SolanaPortfolio
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /solana/portfolio/:sol-address [get]
func SolanaController(c *gin.Context, db *gorm.DB) {
	solAddress := c.Param("sol-address")
	moralisAccessKey := c.GetHeader("x-api-key")

	url := "https://solana-gateway.moralis.io/account/mainnet/" + solAddress + "/portfolio"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-api-key", moralisAccessKey)

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

	var solResponse models.SolanaPortfolio
	if err := json.Unmarshal(body, &solResponse); err != nil {
		log.Fatal(err)
	}

	//TODO: add the data into database

	c.JSON(http.StatusOK, gin.H{
		"message": solResponse,
	})
}
