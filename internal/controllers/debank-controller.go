package controllers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/debank"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@BasePath	/api/v1

// DebankController godoc
//
// @Summary      Fetch Debank Wallet Information
// @Description  Retrieves information for a given Debank address using the BTC.com API.
// @Tags         Debank
// @Accept       json
// @Produce      json
// @Param        addresses  query      array  true  "Debank Address" Format(string)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/debank [get]
func DebankController(c *gin.Context, db *gorm.DB, apiClient providers.APIClient) {
	addresses := c.Query("addresses")
	debankAddresses := strings.Split(addresses, ",")

	if len(debankAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc addresses"))
		return
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	ch := make(chan *models.GlobalWallet, 1) // 1 specifies the buffer size of the channel

	for _, btcAddress := range debankAddresses {
		wg.Add(1)

		go fetchAndSaveDebank(c, db, apiClient, btcAddress, ch, &wg, &mutex)
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

// Fetch and save the data for one address
// TODO: write a common interface in provider which saves the data in database
func fetchAndSaveDebank(c *gin.Context, db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	body, err := apiClient.FetchData(address)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error from api: "+err.Error()))
		return
	}

	// ignore the Address which returns error or empty response
	resp := debank.EvmDebankTotalBalanceApiResponse{}
	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		return
	}

	walletResponse, err := saveDebank(db, address, resp)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error from api"+err.Error()))
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}

// helper to save data
// TODO: write a common interface in provider which saves the data in database
func saveDebank(db *gorm.DB, address string, apiResponse debank.EvmDebankTotalBalanceApiResponse) (*models.GlobalWallet, error) {
	// Begin a new transaction
	tx := db.Begin()

	wallet, err := models.GetOrCreateWallet(tx, address, "Debank")
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Initialize EvmAssetsDebankV1 and set the WalletID
	evmAssetsDebankV1 := models.EvmAssetsDebankV1{
		WalletID:      wallet.WalletID,
		TotalUsdValue: apiResponse.TotalUsdValue,
	}

	// Save EvmAssetsDebankV1
	if err = models.CreateOrUpdateEvmAssetsDebankV1(tx, &evmAssetsDebankV1); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Save Chain
	if err = models.SaveChainDetails(tx, wallet.WalletID, apiResponse.ChainList); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Save token list
	err = models.SaveTokenListByEvmAssetsDebankV1ID(tx, evmAssetsDebankV1.EvmAssetID, apiResponse.TokensList)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Save nft list
	err = models.SaveNFTSListByEvmAssetsDebankV1ID(tx, evmAssetsDebankV1.EvmAssetID, apiResponse.NFTList)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return &models.GlobalWallet{}, err
	}

	// get the data
	walletResponse, err := models.GetGlobalWalletWithEvmDebankInfo(db, address)
	if err != nil {
		return &models.GlobalWallet{}, err
	}

	return walletResponse, nil
}
