package types

import (
	"github.com/0xbase-Corp/portfolio_svc/internal/models"
)

type SolanaApiResponse struct {
	Tokens        []models.Token `json:"tokens"`
	NFTs          []models.NFT   `json:"nfts"`
	NativeBalance struct {
		Lamports string `json:"lamports"`
		Solana   string `json:"solana"`
	} `json:"nativeBalance"`
}
