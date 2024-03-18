package solana

import (
	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/utils"
)

type (
	SolanaApiResponse struct {
		Tokens        []models.Token `json:"tokens"`
		NFTs          []models.NFT   `json:"nfts"`
		NativeBalance struct {
			Lamports string `json:"lamports"`
			Solana   string `json:"solana"`
		} `json:"nativeBalance"`
	}

	SolanaAPI struct{}
)

func (b SolanaAPI) FetchData(address string) ([]byte, error) {
	url := "https://solana-gateway.moralis.io/account/mainnet/" + address
	headers := map[string]string{}

	body, err := utils.CallAPI(url, headers)
	if err != nil {
		return nil, err
	}

	return body, nil
}