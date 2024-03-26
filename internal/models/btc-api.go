package models

import (
	"time"

	"gorm.io/gorm"
)

// BitcoinBtcComV1 represents the bitcoin_btc_com_v1 table.
type (
	BitcoinBtcComV1 struct {
		BtcAssetID  uint      `gorm:"primaryKey;autoIncrement" json:"btc_asset_id"` // Primary key
		WalletID    uint      `gorm:"not null;unique" json:"wallet_id"`             // Foreign key to global_wallets
		BtcUsdPrice float64   `gorm:"type:float" json:"btc_usd_price"`
		UpdatedAt   time.Time `gorm:"" json:"updated_at"`
		CreatedAt   time.Time `gorm:"" json:"created_at"`

		BitcoinAddressInfo *BitcoinAddressInfo `gorm:"foreignKey:BtcAssetID" json:"bitcoin_address_info,omitempty"`
		CoingeckoPriceFeed *CoingeckoPriceFeed `gorm:"-" json:"coingecko_price_feed,omitempty"`
	}

	// BitcoinAddressInfo represents the bitcoin_address_info table.
	BitcoinAddressInfo struct {
		AddressID           uint      `gorm:"primaryKey;autoIncrement" json:"address_id"`
		BtcAssetID          uint      `gorm:"not null" json:"btc_asset_id"`
		Received            float64   `gorm:"type:float" json:"received"`
		Sent                float64   `gorm:"type:float" json:"sent"`
		Balance             float64   `gorm:"type:float" json:"balance"`
		TxCount             int       `gorm:"type:int" json:"tx_count"`
		UnconfirmedTxCount  int       `gorm:"type:int" json:"unconfirmed_tx_count"`
		UnconfirmedReceived float64   `gorm:"type:float" json:"unconfirmed_received"`
		UnconfirmedSent     float64   `gorm:"type:float" json:"unconfirmed_sent"`
		UnspentTxCount      int       `gorm:"type:int" json:"unspent_tx_count"`
		FirstTx             string    `gorm:"type:text" json:"first_tx"`
		LastTx              string    `gorm:"type:text" json:"last_tx"`
		UpdatedAt           time.Time `gorm:"" json:"updated_at"`
		CreatedAt           time.Time `gorm:"" json:"created_at"`
	}
)

// TableName overrides the table name used by GORM to `bitcoin_address_info`
func (BitcoinAddressInfo) TableName() string {
	return "bitcoin_address_info"
}

func (BitcoinBtcComV1) TableName() string {
	return "bitcoin_btc_com_v1"
}

// SaveBitcoinData saves a BitcoinBtcComV1 record and a BitcoinAddressInfo record.
func SaveBitcoinData(tx *gorm.DB, btcAddressInfo *BitcoinAddressInfo, btcComV1 *BitcoinBtcComV1) error {
	// First, handle the BitcoinBtcComV1 record
	existingBtcComV1 := BitcoinBtcComV1{}
	result := tx.Where("wallet_id = ?", btcComV1.WalletID).First(&existingBtcComV1)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// Return any other error encountered during the query.
		return result.Error
	}

	// If the record does not exist, create a new one.
	if result.RowsAffected == 0 {
		if err := tx.Create(btcComV1).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// If the record exists, update it with the new information.
		if err := tx.Model(&existingBtcComV1).Updates(btcComV1).Error; err != nil {
			tx.Rollback()
			return err
		}
		// Use the ID of the existing record for the BitcoinAddressInfo.
		btcComV1.BtcAssetID = existingBtcComV1.BtcAssetID
	}

	// Then, handle the BitcoinAddressInfo record
	if btcAddressInfo != nil {
		// Ensure btcAddressInfo.BtcAssetID is set correctly, either to the newly created btcComV1 ID or an existing one.
		if btcComV1 != nil {
			btcAddressInfo.BtcAssetID = btcComV1.BtcAssetID
		}

		// Check if a BitcoinAddressInfo record already exists for the wallet.
		existingBTCAddressInfo := BitcoinAddressInfo{}
		result := tx.Where("btc_asset_id = ?", btcAddressInfo.BtcAssetID).First(&existingBTCAddressInfo)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			// Return any other error encountered during the query.
			return result.Error
		}

		// If the record does not exist, create a new one.
		if result.RowsAffected == 0 {
			if err := tx.Create(btcAddressInfo).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// If the record exists, update it with the new information.
			if err := tx.Model(&existingBTCAddressInfo).Updates(btcAddressInfo).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}
