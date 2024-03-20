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
	EvmAssetsDebankV1     *EvmAssetsDebankV1     `gorm:"foreignKey:WalletID" json:"evm_assets_debank_v1,omitempty"`
	ChainDetails          *[]ChainDetails        `gorm:"foreignKey:WalletID" json:"chain_details,omitempty"`
}

func (GlobalWallet) TableName() string {
	return "global_wallets"
}

// Get the btc data along with it relations based on btcAddress
func GetGlobalWalletWithBitcoinInfo(tx *gorm.DB, btcAddress string) (*GlobalWallet, error) {
	wallet := GlobalWallet{}

	err := tx.Where("wallet_address = ?", btcAddress).
		Preload("BitcoinBtcComV1").
		Preload("BitcoinBtcComV1.BitcoinAddressInfo").
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// Get the solana data along with it relations based on solAddress
func GetGlobalWalletWithSolanaInfo(tx *gorm.DB, solAddress string) (*GlobalWallet, error) {
	wallet := GlobalWallet{}

	err := tx.Where("wallet_address = ?", solAddress).
		Preload("SolanaAssetsMoralisV1.Tokens").
		Preload("SolanaAssetsMoralisV1.NFTS").
		Preload("SolanaAssetsMoralisV1").
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// Get the debank data along with it relations based on debankAddress
func GetGlobalWalletWithEvmDebankInfo(tx *gorm.DB, debankAddress string) (*GlobalWallet, error) {
	wallet := GlobalWallet{}

	err := tx.Where("wallet_address = ?", debankAddress).
		Preload("ChainDetails").
		Preload("EvmAssetsDebankV1").
		Preload("EvmAssetsDebankV1.TokenList").
		Preload("EvmAssetsDebankV1.NFTList").
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func GetWallet(tx *gorm.DB, walletAddress string) (*GlobalWallet, error) {
	wallet := &GlobalWallet{}

	err := tx.Where("wallet_address = ?", walletAddress).First(&wallet).Error
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func UpdateWalletLastUpdateAt(tx *gorm.DB, wallet *GlobalWallet) error {
	wallet.LastUpdatedAt = time.Now()

	if err := tx.Save(wallet).Error; err != nil {
		return err
	}

	return nil
}

func CreateWallet(tx *gorm.DB, walletAddress, blockchainType string) (*GlobalWallet, error) {
	wallet := &GlobalWallet{
		WalletAddress:  walletAddress,
		BlockchainType: blockchainType,
	}

	if err := tx.Create(&wallet).Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func GetOrCreateWallet(tx *gorm.DB, walletAddress, blockchainType string) (*GlobalWallet, error) {
	wallet, err := GetWallet(tx, walletAddress)
	if err == gorm.ErrRecordNotFound {
		wallet, err = CreateWallet(tx, walletAddress, blockchainType)
	}
	return wallet, err
}
