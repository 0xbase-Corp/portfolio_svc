package controllers

import (
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/types"
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
// @Param        debank-address path string true "Debank Address" Format(string)
// @Param        AccessKey header string true "Debank access key" Format(string)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/debank/{debank-address} [get]
func DebankController(c *gin.Context, db *gorm.DB) {
	debankAddress := c.Param("debank-address")
	debankAccessKey := c.GetHeader("AccessKey")

	if debankAccessKey == "" {
		errors.HandleHttpError(c, errors.NewBadRequestError("Debank access key is missing"))
		return
	}

	totalBalanceApiResponse, err := getDebankTotalBalance(debankAddress, debankAccessKey)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to get Debank total balance: "+err.Error()))
		return
	}

	tokenListApiResponse, err := getDebankTokenList(debankAddress, debankAccessKey)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to get Debank token list: "+err.Error()))
		return
	}

	NFTListApiResponse, err := getDebankNFTList(debankAddress, debankAccessKey)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to get Debank nft list: "+err.Error()))
		return
	}

	tx := db.Begin()

	// Retrieve or save wallet
	wallet, err := models.GetOrCreateWallet(tx, debankAddress, "Debank")
	if err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to retrieve or create wallet: "+err.Error()))
		return
	}

	evmAssetsDebankV1 := models.EvmAssetsDebankV1{
		WalletID:      wallet.WalletID,
		TotalUsdValue: totalBalanceApiResponse.TotalUsdValue,
	}

	// Save EvmAssetsDebankV1
	err = models.CreateOrUpdateEvmAssetsDebankV1(tx, &evmAssetsDebankV1)
	if err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to create/update EvmAssetsDebankV1: "+err.Error()))
		return
	}

	// Save Chain
	err = models.SaveChainDetails(tx, wallet.WalletID, totalBalanceApiResponse.ChainList)
	if err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to create/update Chain details: "+err.Error()))
		return
	}

	// Save token list
	err = models.SaveTokenListByEvmAssetsDebankV1ID(tx, evmAssetsDebankV1.EvmAssetID, tokenListApiResponse)
	if err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to create/update token list: "+err.Error()))
		return
	}

	// Save nft list
	err = models.SaveNFTSListByEvmAssetsDebankV1ID(tx, evmAssetsDebankV1.EvmAssetID, NFTListApiResponse)
	if err != nil {
		tx.Rollback()
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to create/update nft list: "+err.Error()))
		return
	}

	if err := tx.Commit().Error; err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to commit transaction: "+err.Error()))
		return
	}

	walletResponse, err := models.GetGlobalWalletWithEvmDebankInfo(db, debankAddress)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to get wallet data: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, walletResponse)
}

// Get total balance data
func getDebankTotalBalance(debankAddress, debankAccessKey string) (types.EvmDebankTotalBalanceApiResponse, error) {
	chainUrl := "https://pro-openapi.debank.com/v1/user/total_balance?id=" + debankAddress
	headers := map[string]string{
		"Accept":    "application/json",
		"AccessKey": debankAccessKey,
	}

	body, err := utils.CallAPI(chainUrl, headers)
	if err != nil {
		return types.EvmDebankTotalBalanceApiResponse{}, err
	}

	totalBalanceApiResponse := types.EvmDebankTotalBalanceApiResponse{}
	if err := utils.DecodeJSONResponse(body, &totalBalanceApiResponse); err != nil {
		return types.EvmDebankTotalBalanceApiResponse{}, err
	}

	return totalBalanceApiResponse, nil
}

// Get token list data
func getDebankTokenList(debankAddress, debankAccessKey string) ([]*models.TokenList, error) {
	chainUrl := "https://pro-openapi.debank.com/v1/user/all_token_list?id=" + debankAddress
	headers := map[string]string{
		"Accept":    "application/json",
		"AccessKey": debankAccessKey,
	}

	body, err := utils.CallAPI(chainUrl, headers)
	if err != nil {
		return nil, err
	}

	tokens := make([]*models.TokenList, 0)
	if err := utils.DecodeJSONResponse(body, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// Get nft list data
func getDebankNFTList(debankAddress, debankAccessKey string) ([]*models.NFTList, error) {
	chainUrl := "https://pro-openapi.debank.com/v1/user/all_nft_list?id=" + debankAddress
	headers := map[string]string{
		"Accept":    "application/json",
		"AccessKey": debankAccessKey,
	}

	body, err := utils.CallAPI(chainUrl, headers)
	if err != nil {
		return nil, err
	}

	ntfs := make([]*models.NFTList, 0)
	if err := utils.DecodeJSONResponse(body, &ntfs); err != nil {
		return nil, err
	}

	return ntfs, nil
}
