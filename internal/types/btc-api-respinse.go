package types

import (
	"github.com/0xbase-Corp/portfolio_svc/internal/models"
)

type BtcApiResponse struct {
	Data struct {
		Data      models.BitcoinAddressInfo `json:"data"`
		ErrorCode int                       `json:"error_code"`
		ErrNo     int                       `json:"err_no"`
		Message   string                    `json:"message"`
		Status    string                    `json:"status"`
	} `json:"data"`
}
