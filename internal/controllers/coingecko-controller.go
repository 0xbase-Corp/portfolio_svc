package controllers

import (
	"strings"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/coingecko"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
	"gorm.io/gorm"
)

func FetchAndSaveCoingeckoPriceForCrypto(db *gorm.DB, cryptoID, currency string) error {
	priceFeedClient := &coingecko.CoingeckoAPI{}

	body, err := priceFeedClient.FetchData(cryptoID, currency)
	if err != nil {
		return err
	}

	// check the response body if it contains the rate limit error message then ignore to update or create the price feed
	if strings.Contains(string(body), "You've exceeded the Rate Limit") {
		// if rate limit is exceeded, return nil not need to update the price feed
		return nil
	}

	resp := coingecko.CryptoResponse{}

	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		return err
	}

	priceFeed := models.CoingeckoPriceFeed{
		Name:     cryptoID,
		Price:    resp[cryptoID][currency],
		Currency: currency,
	}

	if err := models.UpdateOrCreateCoingeckoPriceFeed(db, &priceFeed); err != nil {
		return err
	}

	return nil

}
