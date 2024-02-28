package types

import "github.com/0xbase-Corp/portfolio_svc/internal/models"

type (
	EvmDebankApiResponse struct {
		TotalUsdValue float64              `json:"total_usd_value"`
		ChainList     *models.ChainDetails `json:"chain_list"`
	}
)
