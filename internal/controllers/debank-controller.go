package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DebankController(c *gin.Context, db *gorm.DB) {
	debankAddress := c.Param("debank-address")

	url := "https://pro-openapi.debank.com/v1/user/total_balance?id=" + debankAddress
	log.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create an HTTP client and execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiResponse := types.EvmDebankApiResponse{}

	// Parse the JSON response into the defined struct.
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Println("Error parsing JSON response:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response: " + err.Error()})
		return
	}
	log.Println("Received response from Debank API")

	var wallet models.GlobalWallet
	err = db.Where("wallet_address = ?", debankAddress).First(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = models.GlobalWallet{
				WalletAddress:  debankAddress,
				BlockchainType: "Debank",
			}
			if err := db.Create(&wallet).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet: " + err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
			return
		}
	}
	log.Println("here is the wallet", wallet)
}
