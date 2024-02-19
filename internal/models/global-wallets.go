package models

import (
	"time"
)

// GlobalWallet represents the global_wallets table.
type GlobalWallet struct {
	WalletID       int       `gorm:"primary_key"`
	PortfolioID    int       `gorm:"not null"`
	WalletAddress  string    `gorm:"type:varchar(255);unique;not null"`
	BlockchainType string    `gorm:"type:varchar(255);not null"`
	APIEndpoint    string    `gorm:"type:text"`
	APIVersion     string    `gorm:"type:varchar(50)"`
	LastUpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
