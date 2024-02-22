package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DebankController(c *gin.Context, db *gorm.DB) {
	log.Println("Debank Controller invoked")

	// Extract the Debank address from the request parameter
	debankAddress := c.Param("debank-address")
	log.Println("for address", debankAddress)

	url := "https://pro-openapi.debank.com/v1/user/total_balance?id=0xD322A0bd6A139cFd359F1EFC540F6cb358d73A16" + debankAddress
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

	var apiResponse struct {
		TotalUsdValue float64 `json:"total_usd_value"`
		ChainList     struct {
			ID               string  `json:"id"`
			CommunityID      int     `json:"community_id"`
			Name             string  `json:"name"`
			NativeTokenID    string  `json:"string"`
			LogoURL          string  `json:"logo_url"`
			WrappedTokenID   string  `json:"wrapped_token_id"`
			IsSupportPreExec bool    `json:"is_support_pre_exec"`
			USDValud         float64 `json:"usd_value"`
		} `json:"chain_list"`
	}

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
