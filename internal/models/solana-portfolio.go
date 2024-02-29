package models

import (
	"time"

	"gorm.io/gorm"
)

// SolanaAssetsMoralisV1 represents the solana_assets_moralis_v1 table.
type (
	SolanaAssetsMoralisV1 struct {
		SolanaAssetID    uint      `gorm:"primaryKey;autoIncrement" json:"solana_asset_id"`
		WalletID         int       `gorm:"not null;unique" json:"wallet_id"`
		Lamports         string    `gorm:"type:varchar(255)" json:"lamports"`
		Solana           string    `gorm:"type:varchar(255)" json:"solana"`
		TotalTokensCount int       `gorm:"type:integer" json:"total_tokens_count"`
		TotalNftsCount   int       `gorm:"type:integer" json:"total_nfts_count"`
		LastUpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"  json:"last_updated_at"`

		// relations use in json responses (optional)
		Tokens *[]Token `gorm:"foreignKey:SolanaAssetID" json:"tokens,omitempty"`
		NFTS   *[]NFT   `gorm:"foreignKey:SolanaAssetID" json:"nfts,omitempty"`
	}

	// Token represents the tokens table.
	Token struct {
		TokenID                int       `gorm:"primary_key" json:"token_id"`
		SolanaAssetID          int       `gorm:"not null" json:"solana_asset_id"`
		AssociatedTokenAddress string    `gorm:"type:varchar(255)" json:"associated_token_address"`
		Mint                   string    `gorm:"type:varchar(255)" json:"mint"`
		AmountRaw              string    `gorm:"type:varchar(255)" json:"amount_raw"`
		Amount                 string    `gorm:"type:varchar(255)" json:"amount"`
		Decimals               string    `gorm:"type:varchar(255)" json:"decimals"`
		Name                   string    `gorm:"type:varchar(255)" json:"name"`
		Symbol                 string    `gorm:"type:varchar(50)" json:"symbol"`
		UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	}

	// NFT represents the nfts table.
	NFT struct {
		NFTID                  int       `gorm:"primary_key" json:"nft_id"`
		SolanaAssetID          int       `gorm:"not null" json:"solana_asset_id"`
		AssociatedTokenAddress string    `gorm:"type:varchar(255)" json:"associated_token_address"`
		Mint                   string    `gorm:"type:varchar(255)" json:"mint"`
		AmountRaw              string    `gorm:"type:varchar(255)" json:"amount_raw"`
		Decimals               string    `gorm:"type:varchar(255)" json:"decimals"`
		Name                   string    `gorm:"type:varchar(255)" json:"name"`
		Symbol                 string    `gorm:"type:varchar(50)" json:"symbol"`
		UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
		UserScore              int       `json:"userscore"`
	}
)

func (SolanaAssetsMoralisV1) TableName() string {
	return "solana_assets_moralis_v1"
}

func (Token) TableName() string {
	return "tokens"
}

func (NFT) TableName() string {
	return "nfts"
}

func SaveSolanaData(tx *gorm.DB, solanaAsset *SolanaAssetsMoralisV1, tokens []Token, nfts []NFT) error {
	// Check if a SolanaAssetsMoralisV1 record already exists for the wallet.
	existingAsset := SolanaAssetsMoralisV1{}
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
