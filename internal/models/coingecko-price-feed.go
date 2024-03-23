package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	// CoingeckoPriceFeed represents a coingecko price feed
	CoingeckoPriceFeed struct {
		ID        int       `gorm:"primaryKey" json:"id"` // database id
		CryptoID  string    `gorm:"type:varchar(255)" json:"crypto_id"`
		Name      string    `gorm:"type:varchar(255)" json:"name"`
		Ticker    string    `gorm:"type:varchar(255)" json:"ticker"`
		UseValue  string    `gorm:"type:varchar(255)" json:"use_value"`
		UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
		CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	}
)

func (CoingeckoPriceFeed) TableName() string {
	return "coingecko_price_feed"
}

// CreateCoingeckoPriceFeed creates a new CoingeckoPriceFeed
func CreateCoingeckoPriceFeed(tx *gorm.DB, coingeckoPriceFeed *CoingeckoPriceFeed) error {
	if result := tx.Create(coingeckoPriceFeed); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetCoingeckoPriceFeedByCryptoID returns a CoingeckoPriceFeed by its crypto_id
func GetCoingeckoPriceFeedByCryptoID(tx *gorm.DB, cryptoID string) (*CoingeckoPriceFeed, error) {
	coingeckoPriceFeed := CoingeckoPriceFeed{}
	err := tx.Where("crypto_id = ?", cryptoID).First(&coingeckoPriceFeed).Error
	if err != nil {
		return nil, err
	}
	return &coingeckoPriceFeed, nil
}

// GetCoingeckoPriceFeedByID returns a CoingeckoPriceFeed by its id
func GetCoingeckoPriceFeedByID(tx *gorm.DB, id int) (*CoingeckoPriceFeed, error) {
	coingeckoPriceFeed := CoingeckoPriceFeed{}
	err := tx.Where("id = ?", id).First(&coingeckoPriceFeed).Error
	if err != nil {
		return nil, err
	}

	return &coingeckoPriceFeed, nil
}
