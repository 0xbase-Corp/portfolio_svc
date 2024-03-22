package controllers

import (
	"net/http"
	"sync"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/bitcoin"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/debank"
	"github.com/0xbase-Corp/portfolio_svc/internal/providers/solana"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	PortfolioAddresses struct {
		BTC []string `json:"btc"`
		Sol []string `json:"sol"`
		EVM []string `json:"evm"`
	}
)

//	@BasePath	/api/v1

// AllPortfolioController godoc
//
// AllPortfolioController defines the route and Swagger annotations for fetching all portfolio information.
// @Summary      Fetch all portfolio information
// @Description  Retrieves information for all portfolios including Bitcoin, Solana, and EVM addresses.
// @Tags         portfolio
// @Accept       json
// @Produce      json
// @Param        addresses body PortfolioAddresses true "Portfolio Addresses"
// @Success      200 {object} map[string][]models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /api/v1/all-portfolio [post]
func AllPortfolioController(c *gin.Context, db *gorm.DB) {
	requestBody := PortfolioAddresses{}

	if err := c.BindJSON(&requestBody); err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error"+err.Error()))
		return
	}

	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	ch := make(chan *models.GlobalWallet, 1)

	bitcoinAPIClient := &bitcoin.BitcoinAPI{}
	solanaAPIClient := &solana.SolanaAPI{}
	debankAPIClient := &debank.DebankAPI{}

	// Handle BTC addresses
	for _, btcAddress := range requestBody.BTC {
		wg.Add(1)
		go fetchAndSaveBtc(c, db, bitcoinAPIClient, btcAddress, ch, &wg, mutex)
	}

	// Handle Solana addresses
	for _, solAddress := range requestBody.Sol {
		wg.Add(1)
		go fetchAndSaveSolana(c, db, solanaAPIClient, solAddress, ch, &wg, mutex)
	}

	// Handle EVM addresses
	for _, evmAddress := range requestBody.EVM {
		wg.Add(1)
		go fetchAndSaveDebank(c, db, debankAPIClient, evmAddress, ch, &wg, mutex)
	}

	// Use a goroutine to close the channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(ch)
	}()

	allResponses := make(map[string][]*models.GlobalWallet)

	// Collect all results from the channel and store them in the map
	for walletResponse := range ch {
		if walletResponse == nil {
			// Handle the case where processing failed for an address
			continue
		}
		allResponses[walletResponse.BlockchainType] = append(allResponses[walletResponse.BlockchainType], walletResponse)
	}

	// Return the map containing responses for each address
	c.JSON(http.StatusOK, allResponses)
}
