package bitcoin

import (
	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
)

type (
	BtcApiResponse struct {
		Data struct {
			Data      models.BitcoinAddressInfo `json:"data"`
			ErrorCode int                       `json:"error_code"`
			ErrNo     int                       `json:"err_no"`
			Message   string                    `json:"message"`
			Status    string                    `json:"status"`
		} `json:"data"`
	}

	BitcoinAPI struct{}
)

func (b *BitcoinAPI) FetchData(address string) ([]byte, error) {
	url := "https://chain.api.btc.com/v3/address/" + address
	headers := map[string]string{}

	body, err := utils.CallAPI(url, headers)
	if err != nil {
		return nil, err
	}

	return body, nil
}
