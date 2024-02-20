package models

import (
	"time"

	"gorm.io/gorm"
)

// BitcoinBtcComV1 represents the bitcoin_btc_com_v1 table.
type BitcoinBtcComV1 struct {
	BtcAssetID  uint      `gorm:"primaryKey;autoIncrement"` // Primary key
	WalletID    uint      `gorm:"not null"`                 // Foreign key to global_wallets
	BtcUsdPrice float64   `gorm:"type:float"`
	UpdatedAt   time.Time `gorm:""`
	CreatedAt   time.Time `gorm:""`
}

// BitcoinAddressInfo represents the bitcoin_address_info table.
type BitcoinAddressInfo struct {
	AddressID           uint      `gorm:"primaryKey;autoIncrement"`
	BtcAssetID          uint      `gorm:"not null"`
	Received            float64   `gorm:"type:float"`
	Sent                float64   `gorm:"type:float"`
	Balance             float64   `gorm:"type:float"`
	TxCount             int       `gorm:"type:int"`
	UnconfirmedTxCount  int       `gorm:"type:int"`
	UnconfirmedReceived float64   `gorm:"type:float"`
	UnconfirmedSent     float64   `gorm:"type:float"`
	UnspentTxCount      int       `gorm:"type:int"`
	FirstTx             string    `gorm:"type:text"`
	LastTx              string    `gorm:"type:text"`
	UpdatedAt           time.Time `gorm:""`
	CreatedAt           time.Time `gorm:""`
}

// TableName overrides the table name used by GORM to `bitcoin_address_info`
func (BitcoinAddressInfo) TableName() string {
	return "bitcoin_address_info"
}

// SaveBitcoinData saves a BitcoinBtcComV1 record and a BitcoinAddressInfo record.
// btcComV1 is optional and can be nil.
func SaveBitcoinData(tx *gorm.DB, btcAddressInfo *BitcoinAddressInfo, btcComV1 *BitcoinBtcComV1) error {
	// First, handle the BitcoinBtcComV1 record
	if btcComV1 != nil {
		var existingBtcComV1 BitcoinBtcComV1
		// Check if a BitcoinBtcComV1 record already exists for the wallet.
		result := tx.Where("wallet_id = ?", btcComV1.WalletID).First(&existingBtcComV1)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			// Return any other error encountered during the query.
			return result.Error
		}

		if result.RowsAffected == 0 {
			// If the record does not exist, create a new one.
			if result := tx.Create(btcComV1); result.Error != nil {
				return result.Error
			}
		} else {
			// If the record exists, update it with the new information.
			if err := tx.Model(&existingBtcComV1).Updates(btcComV1).Error; err != nil {
				return err
			}
			// Use the ID of the existing record for the BitcoinAddressInfo.
			btcComV1.BtcAssetID = existingBtcComV1.BtcAssetID
		}
	}

	// Then, handle the BitcoinAddressInfo record
	if btcAddressInfo != nil {
		// Ensure btcAddressInfo.BtcAssetID is set correctly, either to the newly created btcComV1 ID or an existing one.
		if btcComV1 != nil {
			btcAddressInfo.BtcAssetID = btcComV1.BtcAssetID
		}

		if result := tx.Create(btcAddressInfo); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
