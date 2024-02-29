package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/types"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@BasePath	/api/v1
//
// SolanaController godoc
// @Summary      Fetch Solana portfolio details for a given Solana address
// @Description  Fetch Solana portfolio details, including tokens and NFTs, for a specific Solana address.
// @Tags         solana
// @Accept       json
// @Produce      json
// @Param        sol-address path string true "Solana Address" Format(string)
// @Param        x-api-key header string true "Moralis API Key" Format(string)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/solana/{sol-address} [get]
func SolanaController(c *gin.Context, db *gorm.DB) {
	// Extract the Solana address from the request parameter.
	solAddress := c.Param("sol-address")

	// Extract the Moralis API key from the request header.
	moralisAccessKey := c.GetHeader("x-api-key")

	// Prepare the Moralis API request URL.
	url := "https://solana-gateway.moralis.io/account/mainnet/" + solAddress + "/portfolio"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}

	// Add the Moralis API key to the request header.
	req.Header.Add("x-api-key", moralisAccessKey)

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

	// Define a struct to match the JSON response structure from the Moralis API.
	response := types.SolanaApiResponse{}

	// Parse the JSON response into the defined struct.
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Error parsing JSON response:", err)
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to parse JSON response: "+err.Error()))
		return
	}

	// Check if a wallet with the given Solana address exists in the database.
	wallet := models.GlobalWallet{}
	err = db.Where("wallet_address = ?", solAddress).First(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the wallet doesn't exist, create a new one.
			wallet = models.GlobalWallet{
				WalletAddress:  solAddress,
				BlockchainType: "Solana",
			}
			if err := db.Create(&wallet).Error; err != nil {
				errors.HandleHttpError(c, errors.NewBadRequestError("Failed to create wallet: "+err.Error()))
				return
			}
		} else {
			errors.HandleHttpError(c, errors.NewBadRequestError("Database query error: "+err.Error()))
			return
		}
	}

	// Prepare the SolanaAssetsMoralisV1 object with the response data.
	solanaAsset := models.SolanaAssetsMoralisV1{
		WalletID:         wallet.WalletID, // Assuming WalletID is the correct field name
		Lamports:         response.NativeBalance.Lamports,
		Solana:           response.NativeBalance.Solana,
		TotalTokensCount: len(response.Tokens),
		TotalNftsCount:   len(response.NFTs),
		LastUpdatedAt:    time.Now(),
	}

	// Start a new database transaction.
	tx := db.Begin()

	// Attempt to save the Solana asset data along with the associated tokens and NFTs.
	if err := models.SaveSolanaData(tx, &solanaAsset, response.Tokens, response.NFTs); err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to save data to the database: "+err.Error()))
		return
	}

	// Commit the transaction if everything is successful.
	if err := tx.Commit().Error; err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to commit transaction: "+err.Error()))
		return
	}

	walletResponse, err := models.GetGlobalWalletWithSolanaInfo(db, solAddress)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to get wallet data: "+err.Error()))
		return
	}

	// Send a success response.
	c.JSON(http.StatusOK, walletResponse)
}

//	@BasePath	/api/v1
//
// GetSolanaController godoc
// @Summary      Get Solana portfolio for a wallet
// @Description  Retrieve Solana portfolio details, including tokens and NFTs, for a specific wallet.
// @Tags         solana
// @Accept       json
// @Produce      json
// @Param        wallet_id path int true "Wallet ID" Format(int)
// @Param        offset query int false "Pagination offset" Format(int)
// @Param        limit query int false "Pagination limit" Format(int)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/solana-wallet/{wallet_id} [get]
func GetSolanaController(c *gin.Context, db *gorm.DB) {
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

	/**
	explain preload query

	SELECT *
	FROM global_wallets
	LEFT JOIN solana_assets_moralis_v1 ON global_wallets.wallet_id = solana_assets_moralis_v1.wallet_id
	LEFT JOIN tokens ON solana_assets_moralis_v1.solana_asset_id = tokens.solana_asset_id
	LEFT JOIN nfts ON solana_assets_moralis_v1.solana_asset_id = nfts.solana_asset_id
	WHERE global_wallets.wallet_id = ?
	LIMIT limit;
	*/

	// Define a reusable function to apply offset and limit for preload
	prlimit := func(query *gorm.DB) *gorm.DB {
		return query.Offset((page - 1) * limit).Limit(limit)
	}

	err = db.
		Preload("SolanaAssetsMoralisV1.Tokens", prlimit).
		Preload("SolanaAssetsMoralisV1.NFTS", prlimit).
		Preload("SolanaAssetsMoralisV1").
		Where("wallet_id = ?", walletID).
		First(&wallet).Error

	if err != nil {
		errors.HandleHttpError(c, errors.NewNotFoundError("wallet not found"))
		return
	}

	c.JSON(http.StatusOK, wallet)
}
