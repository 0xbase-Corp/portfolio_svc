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

// SaveBitcoinData saves a BitcoinBtcComV1 record and a BitcoinAddressInfo record.
// btcComV1 is optional and can be nil.
func SaveBitcoinData(tx *gorm.DB, btcAddressInfo *BitcoinAddressInfo, btcComV1 *BitcoinBtcComV1) error {
	// Check if btcComV1 is not nil before saving
	if btcComV1 != nil {
		if result := tx.Create(btcComV1); result.Error != nil {
			return result.Error
		}
	}

	// Check if btcAddressInfo is not nil before saving
	if btcAddressInfo != nil {
		if result := tx.Create(btcAddressInfo); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
