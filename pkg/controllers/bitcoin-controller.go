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

// BitcoinController handles requests for Bitcoin wallet information.
func BitcoinController(c *gin.Context, db *gorm.DB) {
	log.Println("BitcoinController invoked")
	// Extract the Bitcoin address from the request parameter.
	btcAddress := c.Param("btc-address")
	log.Println("for address:", btcAddress)

	// Prepare the BTC API request URL.
	url := "https://chain.api.btc.com/v3/address/" + btcAddress
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

	// Define a struct to match the JSON response structure from the BTC API.
	var apiResponse struct {
		Data struct {
			Data      models.BitcoinAddressInfo `json:"data"`
			ErrorCode int                       `json:"error_code"`
			ErrNo     int                       `json:"err_no"`
			Message   string                    `json:"message"`
			Status    string                    `json:"status"`
		} `json:"data"`
	}

	// Parse the JSON response into the defined struct.
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Println("Failed to parse JSON response: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response: " + err.Error()})
		return
	}
	log.Println("Received response from BTC API")

	// Check if a wallet with the given Bitcoin address exists in the global_wallets table
	var wallet models.GlobalWallet
	err = db.Where("wallet_address = ?", btcAddress).First(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = models.GlobalWallet{
				WalletAddress:  btcAddress,
				BlockchainType: "Bitcoin",
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

	// Begin a new transaction
	tx := db.Begin()
	// Assuming wallet is the GlobalWallet record found or created
	walletID := wallet.WalletID

	// Initialize btcComV1 and set the WalletID
	var btcComV1 models.BitcoinBtcComV1
	btcComV1.WalletID = uint(walletID)

	// Save btcComV1 to the database to get a valid BtcAssetID
	if err := tx.Create(&btcComV1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save btcComV1 record: " + err.Error()})
		return
	}

	// Now that btcComV1 is saved, it has a valid BtcAssetID. Set this ID for BitcoinAddressInfo
	apiResponse.Data.Data.BtcAssetID = btcComV1.BtcAssetID

	// Proceed to call SaveBitcoinData with the updated BitcoinAddressInfo and btcComV1
	if err := models.SaveBitcoinData(tx, &apiResponse.Data.Data, &btcComV1); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Bitcoin address info: " + err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	log.Println("Data saved successfully for address:", btcAddress)
	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})
}
