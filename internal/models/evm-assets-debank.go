package models

import (
	"time"
)

type (
	EvmAssetsDebankV1 struct {
		EvmAssetId    uint      `gorm:"primaryKey;autoIncrement" json:"evm_asset_id"` // Primary key
		WalletID      uint      `gorm:"not null;unique" json:"wallet_id"`             // Foreign key to global_wallets
		TotalUsdValue float64   `gorm:"type:float" json:"total_usd_value"`
		ChainListJson string    `gorm:"type:text" json:"chain_list_json"`
		UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	}

	ChainDetails struct {
		ChainID        string    `gorm:"primaryKey" json:"chain_id"`
		WalletID       uint      `gorm:"not null;unique" json:"wallet_id"`
		CommunityID    string    `gorm:"type:varchar(255)" json:"community_id"`
		Name           string    `gorm:"type:varchar(255)" json:"name"`
		LogoURL        string    `gorm:"type:varchar(255)" json:"logo_url"`
		NativeTokenID  string    `gorm:"type:varchar(255)" json:"native_token_id"`
		WrappedTokenID string    `gorm:"type:varchar(255)" json:"wrapped_token_id"`
		USDValue       float64   `gorm:"type:decimal(10,2)" json:"usd_value"`
		UpdatedAt      time.Time `json:"updated_at" gorm:"default:current_timestamp"`
		CreatedAt      time.Time `json:"created_at" gorm:"default:current_timestamp"`
	}

	TokenList struct {
		TokenID      string    `gorm:"primaryKey" json:"token_id"`
		EvmAssetID   uint      `gorm:"not null" json:"evm_asset_id"`
		ContractID   string    `gorm:"type:varchar(255)" json:"contract_id"`
		InnerID      string    `gorm:"type:varchar(255)" json:"inner_id"`
		Chain        string    `gorm:"type:varchar(255)" json:"chain"`
		Name         string    `gorm:"type:varchar(255)" json:"name"`
		Description  string    `gorm:"type:text" json:"description"`
		ContentType  string    `gorm:"type:varchar(255)" json:"content_type"`
		Content      string    `gorm:"type:text" json:"content"`
		DetailURL    string    `gorm:"type:text" json:"detail_url"`
		ContractName string    `gorm:"type:varchar(255)" json:"contract_name"`
		IsERC1155    bool      `gorm:"type:boolean" json:"is_erc1155"`
		Amount       float64   `gorm:"type:float" json:"amount"`
		ProtocolJSON string    `gorm:"type:text" json:"protocol_json"`
		PayTokenJSON string    `gorm:"type:text" json:"pay_token_json"`
		CollectionID string    `gorm:"type:varchar(255)" json:"collection_id"`
		UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	}

	NFTList struct {
		AttributeID uint      `gorm:"primaryKey" json:"attribute_id"`
		EvmAssetID  int       `gorm:"not null" json:"evm_asset_id"`
		NFTID       string    `gorm:"type:varchar(255)" json:"nft_id"`
		Key         string    `gorm:"type:varchar(255)" json:"key"`
		TraitType   string    `gorm:"type:varchar(255)" json:"trait_type"`
		Value       string    `gorm:"type:text" json:"value"`
		UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
		CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	}
)

func (EvmAssetsDebankV1) TableName() string {
	return "evm_assets_debank_v1"
}

func (ChainDetails) TableName() string {
	return "chain_details"
}

func (TokenList) TableName() string {
	return "token_list"
}

func (NFTList) TableName() string {
	return "nft_list"
}
