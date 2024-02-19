package models

import (
	"time"

	"gorm.io/gorm"
)

// SolanaAssetsMoralisV1 represents the solana_assets_moralis_v1 table.
type SolanaAssetsMoralisV1 struct {
	SolanaAssetID    uint      `gorm:"primaryKey;autoIncrement"`
	WalletID         int       `gorm:"not null;unique"`
	Lamports         string    `gorm:"type:varchar(255)"`
	Solana           string    `gorm:"type:varchar(255)"`
	TotalTokensCount int       `gorm:"type:integer"`
	TotalNftsCount   int       `gorm:"type:integer"`
	LastUpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Token represents the tokens table.
type Token struct {
	TokenID                int       `gorm:"primary_key"`
	SolanaAssetID          int       `gorm:"not null"`
	AssociatedTokenAddress string    `gorm:"type:varchar(255)"`
	Mint                   string    `gorm:"type:varchar(255):"`
	AmountRaw              string    `gorm:"type:varchar(255)"`
	Amount                 string    `gorm:"type:varchar(255)"`
	Decimals               string    `gorm:"type:varchar(255)"`
	Name                   string    `gorm:"type:varchar(255)"`
	Symbol                 string    `gorm:"type:varchar(50)"`
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// NFT represents the nfts table.
type NFT struct {
	NFTID                  int       `gorm:"primary_key"`
	SolanaAssetID          int       `gorm:"not null"`
	AssociatedTokenAddress string    `gorm:"type:varchar(255)"`
	Mint                   string    `gorm:"type:varchar(255)"`
	AmountRaw              string    `gorm:"type:varchar(255)"`
	Decimals               string    `gorm:"type:varchar(255)"`
	Name                   string    `gorm:"type:varchar(255)"`
	Symbol                 string    `gorm:"type:varchar(50)"`
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func SaveSolanaData(tx *gorm.DB, solanaAsset *SolanaAssetsMoralisV1, tokens []Token, nfts []NFT) error {
	// Check if a SolanaAssetsMoralisV1 record already exists for the wallet.
	var existingAsset SolanaAssetsMoralisV1
	result := tx.Where("wallet_id = ?", solanaAsset.WalletID).First(&existingAsset)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// Return any other error encountered during the query.
		return result.Error
	}

	if result.RowsAffected == 0 {
		// If the record does not exist, create a new one.
		if result := tx.Create(solanaAsset); result.Error != nil {
			return result.Error
		}
	} else {
		// If the record exists, update it with the new information.
		if err := tx.Model(&existingAsset).Updates(solanaAsset).Error; err != nil {
			return err
		}
		// Use the ID of the existing record for associated tokens and NFTs.
		solanaAsset.SolanaAssetID = existingAsset.SolanaAssetID
	}

	// Set SolanaAssetID for each Token and NFT.
	for i := range tokens {
		tokens[i].SolanaAssetID = int(solanaAsset.SolanaAssetID)
	}
	for i := range nfts {
		nfts[i].SolanaAssetID = int(solanaAsset.SolanaAssetID)
	}

	// Save each Token
	for _, token := range tokens {
		if result := tx.Create(&token); result.Error != nil {
			return result.Error
		}
	}

	// Save each NFT
	for _, nft := range nfts {
		if result := tx.Create(&nft); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
