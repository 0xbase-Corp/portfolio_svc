package controllers

import (
	"strings"
	"time"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/coingecko"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
	"gorm.io/gorm"
)

const maxRetries = 3

// TODO: handle rate limit errors
// TODO: add timer to call this service
func FetchAndSaveCoingeckoPriceForCrypto(db *gorm.DB, cryptoID, currency string) error {
	priceFeedClient := &coingecko.CoingeckoAPI{}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		body, err := priceFeedClient.FetchData(cryptoID, currency)
		if err != nil {
			return err
		}

		bodyStr := string(body)

		if strings.Contains(bodyStr, "You've exceeded the Rate Limit") {
			time.Sleep(time.Second * time.Duration(2^attempt)) // Exponential backoff
			continue
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

	return nil
}
