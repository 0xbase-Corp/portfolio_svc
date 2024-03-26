package coingecko

import (
	"fmt"

	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
)

type (
	CryptoResponse map[string]map[string]float64

	CoingeckoAPI struct{}
)

func (c *CoingeckoAPI) FetchData(cryptoID, currency string) ([]byte, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", cryptoID, currency)

	headers := map[string]string{}

	body, err := utils.CallAPI(url, headers)

	if err != nil {
		return nil, err
	}

	return body, nil
}
