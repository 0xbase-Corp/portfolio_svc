package models

import (
	"time"

	"gorm.io/gorm"
)

// GlobalWallet represents the global_wallets table.
type GlobalWallet struct {
	WalletID       int       `gorm:"primary_key" json:"wallet_id"`
	PortfolioID    int       `gorm:"not null" json:"portfolio_id"`
	WalletAddress  string    `gorm:"type:varchar(255);unique;not null" json:"wallet_address"`
	BlockchainType string    `gorm:"type:varchar(255);not null" json:"blockchain_type"`
	APIEndpoint    string    `gorm:"type:text" json:"api_endpoint"`
	APIVersion     string    `gorm:"type:varchar(50)" json:"api_version"`
	LastUpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_updated_at"`

	// relations use in json responses (optional)
	SolanaAssetsMoralisV1 *SolanaAssetsMoralisV1 `gorm:"foreignKey:WalletID" json:"solana_assets_moralis_v1,omitempty"`
	BitcoinBtcComV1       *BitcoinBtcComV1       `gorm:"foreignKey:WalletID" json:"bitcoin_btc_com_v1,omitempty"`
}

func (GlobalWallet) TableName() string {
	return "global_wallets"
}

// Get the btc data along with it relations based on btcAddress
func GetGlobalWalletWithBitcoinInfo(db *gorm.DB, btcAddress string) (*GlobalWallet, error) {
	wallet := GlobalWallet{}

	err := db.Where("wallet_address = ?", btcAddress).
		Preload("BitcoinBtcComV1").
		Preload("BitcoinBtcComV1.BitcoinAddressInfo").
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func GetGlobalWalletWithSolanaInfo(db *gorm.DB, solAddress string) (*GlobalWallet, error) {
	wallet := GlobalWallet{}

	err := db.Where("wallet_address = ?", solAddress).
		Preload("SolanaAssetsMoralisV1.Tokens").
		Preload("SolanaAssetsMoralisV1.NFTS").
		Preload("SolanaAssetsMoralisV1").
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
