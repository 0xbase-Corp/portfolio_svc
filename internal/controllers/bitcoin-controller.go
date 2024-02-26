package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/types"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@BasePath	/api/v1

// BitcoinController godoc
//
// @Summary      Fetch Bitcoin Wallet Information
// @Description  Retrieves information for a given Bitcoin address using the BTC.com API.
// @Tags         bitcoin
// @Accept       json
// @Produce      json
// @Param        btc-address  path      string  true  "Bitcoin Address"
// @Success      200 {object} types.OKResponse
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/btc/{btc-address} [get]
func BitcoinController(c *gin.Context, db *gorm.DB) {
	btcAddress := c.Param("btc-address")

	if btcAddress == "" {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc btc address"))
		return
	}

	// Prepare the BTC API request URL.
	url := "https://chain.api.btc.com/v3/address/" + btcAddress
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}
	// Create an HTTP client and execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}
	defer resp.Body.Close()
	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
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
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to parse JSON response: "+err.Error()))
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
				errors.HandleHttpError(c, errors.NewInternalServerError("Failed to create wallet: "+err.Error()))
				return
			}
		} else {
			errors.HandleHttpError(c, errors.NewInternalServerError("Database query error: "+err.Error()))
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
		errors.HandleHttpError(c, errors.NewInternalServerError("Failed to save btcComV1 record: "+err.Error()))
		return
	}

	// Now that btcComV1 is saved, it has a valid BtcAssetID. Set this ID for BitcoinAddressInfo
	apiResponse.Data.Data.BtcAssetID = btcComV1.BtcAssetID

	// Proceed to call SaveBitcoinData with the updated BitcoinAddressInfo and btcComV1
	if err := models.SaveBitcoinData(tx, &apiResponse.Data.Data, &btcComV1); err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewInternalServerError("Failed to save Bitcoin address info: "+err.Error()))
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		errors.HandleHttpError(c, errors.NewInternalServerError("Failed to commit transaction: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, types.OKResponse{
		Message: "Data saved successfully",
	})
}

//	@BasePath	/api/v1
//
// GetBtcDataController godoc
// @Summary      Get BTC portfolio for a wallet
// @Description  Retrieve BTC portfolio details, including BitcoinAddressInfo, for a specific wallet.
// @Tags         bitcoin
// @Accept       json
// @Produce      json
// @Param        wallet_id path int true "Wallet ID" Format(int)
// @Param        offset query int false "Pagination offset" Format(int)
// @Param        limit query int false "Pagination limit" Format(int)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/btc-wallet/{wallet_id} [get]
func GetBtcDataController(c *gin.Context, db *gorm.DB) {
	wallet := models.GlobalWallet{}
	walletID, err := strconv.Atoi(c.Param("wallet-id"))

	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("invalid wallet id"))
		return
	}

	// Parse optional query parameters
	page, err := strconv.Atoi(c.DefaultQuery("offset", "1"))
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("invalid offset"))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("invalid limit"))
		return
	}

	// Define a reusable function to apply offset and limit for preload on BitcoinAddressInfo
	prlimit := func(query *gorm.DB) *gorm.DB {
		return query.Offset((page - 1) * limit).Limit(limit)
	}

	err = db.
		Preload("BitcoinBtcComV1.BitcoinAddressInfo", prlimit).
		Preload("BitcoinBtcComV1").
		Where("wallet_id = ?", walletID).
		First(&wallet).Error

	if err != nil {
		errors.HandleHttpError(c, errors.NewNotFoundError("wallet not found"))
		return
	}

	c.JSON(http.StatusOK, wallet)
}
