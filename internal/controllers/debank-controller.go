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
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to get Debank token list: "+err.Error()))
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

	// Save EvmAssetsDebankV1
	evmAssetsDebankV1 := models.EvmAssetsDebankV1{
		WalletID:      wallet.WalletID,
		TotalUsdValue: totalBalanceApiResponse.TotalUsdValue,
	}

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

	// create or update token list

	if err := tx.Commit().Error; err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("Failed to commit transaction: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokenListApiResponse,
		"nfts":   NFTListApiResponse,
	})
}

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
