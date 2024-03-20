package debank

import (
	"encoding/json"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
	"github.com/0xbase-Corp/portfolio_svc/shared/configs"
)

type (
	EvmDebankTotalBalanceApiResponse struct {
		TotalUsdValue float64                `json:"total_usd_value"`
		ChainList     []*models.ChainDetails `json:"chain_list"`
		TokensList    []*models.TokenList    `json:"token_list"`
		NFTList       []*models.NFTList      `json:"nft_list"`
	}

	DebankAPI struct{}
)

func (d *DebankAPI) FetchData(address string) ([]byte, error) {
	headers := map[string]string{
		"Accept":    "application/json",
		"AccessKey": configs.EnvConfigVars.GetDebankAccessKeyHeader(),
	}

	resp := EvmDebankTotalBalanceApiResponse{}
	if err := d.fetch("https://pro-openapi.debank.com/v1/user/total_balance?id="+address, headers, &resp); err != nil {
		return nil, err
	}

	if err := d.fetch("https://pro-openapi.debank.com/v1/user/all_token_list?id="+address, headers, &resp.TokensList); err != nil {
		return nil, err
	}

	if err := d.fetch("https://pro-openapi.debank.com/v1/user/all_nft_list?id="+address, headers, &resp.NFTList); err != nil {
		return nil, err
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (d *DebankAPI) fetch(url string, headers map[string]string, data interface{}) error {
	body, err := utils.CallAPI(url, headers)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}
