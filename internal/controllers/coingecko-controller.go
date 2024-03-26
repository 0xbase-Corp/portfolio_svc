package controllers

import (
	"time"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/coingecko"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
	"gorm.io/gorm"
)

func HandleCoingeckoPrice(tx *gorm.DB, cryptoID, currency string) error {
	fetched, _ := models.GetCoingeckoPriceFeedByName(tx, cryptoID)

	if fetched != nil {
		if timeSinceLastUpdate := time.Since(fetched.UpdatedAt); timeSinceLastUpdate.Minutes() > 2 {
			if err := fetchAndSaveCoingeckoPriceForCrypto(tx, cryptoID, currency); err != nil {
				return err
			}
		}
	} else {
		if err := fetchAndSaveCoingeckoPriceForCrypto(tx, cryptoID, currency); err != nil {
			return err
		}
	}

	return nil
}

func fetchAndSaveCoingeckoPriceForCrypto(db *gorm.DB, cryptoID, currency string) error {
	priceFeedClient := &coingecko.CoingeckoAPI{}

	body, err := priceFeedClient.FetchData(cryptoID, currency)
	if err != nil {
		return err
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
