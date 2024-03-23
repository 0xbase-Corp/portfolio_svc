package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/solana"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
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
// @Param        addresses  query      array  true  "Solana Addresses" Format(string)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/solana [get]
func SolanaController(c *gin.Context, db *gorm.DB, apiClient providers.APIClient) {
	addresses := c.Query("addresses")
	solanaAddresses := strings.Split(addresses, ",")

	if len(solanaAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc addresses"))
		return
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	ch := make(chan *models.GlobalWallet, 1) // 1 specifies the buffer size of the channel

	for _, solAddress := range solanaAddresses {
		wg.Add(1)

		go fetchAndSaveSolana(c, db, apiClient, solAddress, ch, &wg, &mutex)
	}

	// Use a goroutine to close the channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect all results from the channel
	responses := make([]*models.GlobalWallet, 0)
	for walletResponse := range ch {
		if walletResponse == nil {
			// Handle the case where processing failed for an address
			continue
		}
		responses = append(responses, walletResponse)
	}

	c.JSON(http.StatusOK, responses)
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

// helper to save data
// TODO: write a common interface in provider which saves the data in database
func saveSolana(db *gorm.DB, solanaAddress string, apiResponse solana.SolanaApiResponse) (*models.GlobalWallet, error) {
	// Begin a new transaction
	tx := db.Begin()

	wallet, err := models.GetOrCreateWallet(tx, solanaAddress, "Solana")
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Assuming wallet is the GlobalWallet record found or created
	walletID := wallet.WalletID

	// Initialize solanaAsset and set the WalletID
	solanaAsset := models.SolanaAssetsMoralisV1{}
	solanaAsset.WalletID = walletID

	// Attempt to save the Solana asset data along with the associated tokens and NFTs.
	if err := models.SaveSolanaData(tx, &solanaAsset, apiResponse.Tokens, apiResponse.NFTs); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return &models.GlobalWallet{}, err
	}

	// get the data
	walletResponse, err := models.GetGlobalWalletWithSolanaInfo(db, solanaAddress)
	if err != nil {
		return &models.GlobalWallet{}, err
	}

	return walletResponse, nil
}

// Fetch and save the data for one address
// TODO: write a common interface in provider which saves the data in database
func fetchAndSaveSolana(c *gin.Context, db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	body, err := apiClient.FetchData(address)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error from api: "+err.Error()))
		return
	}

	// ignore the Address which returns error or empty response
	resp := solana.SolanaApiResponse{}
	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		return
	}

	walletResponse, err := saveSolana(db, address, resp)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error from api"+err.Error()))
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}
