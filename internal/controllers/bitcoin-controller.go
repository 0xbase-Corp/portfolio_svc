package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/bitcoin"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
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
// @Param        addresses  query      array  true  "Bitcoin Addresses"
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/btc [get]
func BitcoinController(c *gin.Context, db *gorm.DB, apiClient providers.APIClient) {
	addresses := c.Query("addresses")
	btcAddresses := strings.Split(addresses, ",")

	if len(btcAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc addresses"))
		return
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	ch := make(chan *models.GlobalWallet, 1) // 1 specifies the buffer size of the channel

	for _, btcAddress := range btcAddresses {
		wg.Add(1)

		go process(c, db, apiClient, btcAddress, ch, &wg, &mutex)
	}

	// Use a goroutine to close the channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect all results from the channel
	allResponses := make([]*models.GlobalWallet, 0)
	for walletResponse := range ch {
		if walletResponse == nil {
			// Handle the case where processing failed for an address
			continue
		}
		allResponses = append(allResponses, walletResponse)
	}

	c.JSON(http.StatusOK, allResponses)
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

// helper to save data
func save(db *gorm.DB, btcAddress string, apiResponse bitcoin.BtcApiResponse) (*models.GlobalWallet, error) {
	// Begin a new transaction
	tx := db.Begin()

	wallet, err := models.GetOrCreateWallet(tx, btcAddress, "Bitcoin")
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Assuming wallet is the GlobalWallet record found or created
	walletID := wallet.WalletID

	// Initialize btcComV1 and set the WalletID
	btcComV1 := models.BitcoinBtcComV1{}
	btcComV1.WalletID = uint(walletID)

	// Proceed to call SaveBitcoinData with the updated BitcoinAddressInfo and btcComV1
	if err := models.SaveBitcoinData(tx, &apiResponse.Data.Data, &btcComV1); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return &models.GlobalWallet{}, err
	}

	// get the data
	walletResponse, err := models.GetGlobalWalletWithBitcoinInfo(db, btcAddress)
	if err != nil {
		return &models.GlobalWallet{}, err
	}

	return walletResponse, nil
}

// Fetch and save the data for one address
func process(c *gin.Context, db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	body, err := apiClient.FetchData(address)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error from api: "+err.Error()))
		return
	}

	// resp := bitcoin.BtcApiResponse{}
	// if err := utils.DecodeJSONResponse(body, &resp); err != nil {
	// 	errors.HandleHttpError(c, errors.NewBadRequestError("Failed to parse JSON response: "+err.Error()))
	// 	return
	// }

	// ignore the Address which returns error or empty response
	resp := bitcoin.BtcApiResponse{}
	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		return
	}

	walletResponse, err := save(db, address, resp)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error from api"+err.Error()))
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}
