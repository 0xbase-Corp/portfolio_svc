package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	EvmAssetsDebankV1 struct {
		EvmAssetID    uint      `gorm:"primaryKey;autoIncrement" json:"evm_asset_id"` // Primary key
		WalletID      int       `gorm:"not null;unique" json:"wallet_id"`             // Foreign key to global_wallets
		TotalUsdValue float64   `gorm:"type:float" json:"total_usd_value"`
		ChainListJson string    `gorm:"type:text" json:"chain_list_json"`
		UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	}

	ChainDetails struct {
		ChainID        int       `gorm:"primaryKey" json:"chian_id"`
		ID             string    `gorm:"type:varchar(255)" json:"id"`
		WalletID       int       `gorm:"not null" json:"wallet_id"`
		CommunityID    uint64    `gorm:"type:integer" json:"community_id"`
		Name           string    `gorm:"type:varchar(255)" json:"name"`
		LogoURL        string    `gorm:"type:varchar(255)" json:"logo_url"`
		NativeTokenID  string    `gorm:"type:varchar(255)" json:"native_token_id"`
		WrappedTokenID string    `gorm:"type:varchar(255)" json:"wrapped_token_id"`
		USDValue       float64   `gorm:"type:decimal(10,2)" json:"usd_value"`
		UpdatedAt      time.Time `json:"updated_at" gorm:"default:current_timestamp"`
		CreatedAt      time.Time `json:"created_at" gorm:"default:current_timestamp"`
	}

	TokenList struct {
		TokenID      int       `gorm:"primaryKey" json:"token_id"`  // database id
		ID           string    `gorm:"type:varchar(255)" json:"id"` // token id
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
		NFTID        string    `gorm:"primaryKey" json:"nft_id"`    // database id
		ID           string    `gorm:"type:varchar(255)" json:"id"` // nft id
		EvmAssetID   int       `gorm:"not null" json:"evm_asset_id"`
		ContractID   string    `gorm:"type:varchar(255)" json:"contract_id"`
		InnerID      string    `gorm:"type:varchar(255)" json:"inner_id"`
		Chain        string    `gorm:"type:varchar(255)" json:"chain"`
		Name         string    `gorm:"type:varchar(255)" json:"name"`
		Description  string    `gorm:"type:text" json:"description"`
		ContentType  string    `gorm:"type:varchar(255)" json:"content_type"`
		Content      string    `gorm:"type:text" json:"content"`
		ThumbnailURL string    `gorm:"type:varchar(255)" json:"thumbnail_url"`
		TotalSupply  int       `gorm:"type:integer" json:"total_supply"`
		DetailURL    string    `gorm:"type:varchar(255)" json:"detail_url"`
		CollectionID string    `gorm:"type:varchar(255)" json:"collection_id"`
		ContractName string    `gorm:"type:varchar(255)" json:"contract_name"`
		IsERC1155    bool      `gorm:"type:boolean" json:"is_erc1155"`
		Amount       int       `gorm:"type:integer" json:"amount"`
		USDPrice     float64   `gorm:"type:float" json:"usd_price"`
		Attributes   any       `gorm:"type:json" json:"attributes"`
		PayToken     any       `gorm:"type:json" json:"pay_token"`
		UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
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

func CreateEvmAssetsDebankV1(tx *gorm.DB, evmAssetsDebankV1 *EvmAssetsDebankV1) error {
	if result := tx.Create(evmAssetsDebankV1); result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateEvmAssetsDebankV1(tx *gorm.DB, existingAsset *EvmAssetsDebankV1, evmAssetsDebankV1 *EvmAssetsDebankV1) error {
	if err := tx.Model(existingAsset).Updates(evmAssetsDebankV1).Error; err != nil {
		return err
	}
	return nil
}

func GetEvmAssetsDebankV1ByWalletID(tx *gorm.DB, walletID int) (EvmAssetsDebankV1, error) {
	var existingAsset EvmAssetsDebankV1
	result := tx.Where("wallet_id = ?", walletID).First(&existingAsset)
	return existingAsset, result.Error
}

func CreateOrUpdateEvmAssetsDebankV1(tx *gorm.DB, evmAssetsDebankV1 *EvmAssetsDebankV1) error {
	existingAsset, err := GetEvmAssetsDebankV1ByWalletID(tx, evmAssetsDebankV1.WalletID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		if err := CreateEvmAssetsDebankV1(tx, evmAssetsDebankV1); err != nil {
			return err
		}
	} else {
		if err := UpdateEvmAssetsDebankV1(tx, &existingAsset, evmAssetsDebankV1); err != nil {
			return err
		}
		// Use the ID of the existing record for associated tokens and NFTs.
		evmAssetsDebankV1.EvmAssetID = existingAsset.EvmAssetID
	}

	return nil
}

func FindChainsByWalletID(tx *gorm.DB, walletID int) (*[]ChainDetails, error) {
	var existingChains *[]ChainDetails
	if err := tx.Where("wallet_id = ?", walletID).Find(&existingChains).Error; err != nil {
		return nil, err
	}
	return existingChains, nil
}

// SaveChainDetails finds a chain by WalletID, if exists, updates it; otherwise, creates a new chain
func SaveChainDetails(tx *gorm.DB, walletID int, chainDetails []*ChainDetails) error {
	// Find existing chains by wallet ID
	existingChains, err := FindChainsByWalletID(tx, walletID)
	if err != nil {
		return err
	}

	// Create a map to store existing chains by ID for efficient lookup
	existingChainByID := make(map[string]ChainDetails)
	for _, existingChain := range *existingChains {
		existingChainByID[existingChain.ID] = existingChain
	}

	// Iterate through the provided chain details
	for _, chainDetail := range chainDetails {
		// Check if the chain exists
		existingChain, exists := existingChainByID[chainDetail.ID]

		if !exists {
			// If the chain does not exist, create a new one
			chainDetail.WalletID = walletID
			if err := tx.Create(&chainDetail).Error; err != nil {
				return err
			}
		} else {
			// If the chain exists, update its details
			chainDetail.WalletID = walletID
			chainDetail.ChainID = existingChain.ChainID
			if err := tx.Model(&existingChain).Updates(chainDetail).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
