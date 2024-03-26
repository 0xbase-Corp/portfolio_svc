package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	// CoingeckoPriceFeed represents a coingecko price feed
	CoingeckoPriceFeed struct {
		ID        int       `gorm:"primaryKey" json:"id"` // database id
		Name      string    `gorm:"type:varchar(255)" json:"name"`
		Price     float64   `gorm:"type:float" json:"price"`
		Currency  string    `gorm:"type:varchar(50)" json:"currency"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
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

// UpdateOrCreateCoingeckoPriceFeed updates or create a CoingeckoPriceFeed
func UpdateOrCreateCoingeckoPriceFeed(tx *gorm.DB, coingeckoPriceFeed *CoingeckoPriceFeed) error {
	existingCoingeckoPriceFeed, _ := GetCoingeckoPriceFeedByName(tx, coingeckoPriceFeed.Name)

	if existingCoingeckoPriceFeed == nil {
		if err := CreateCoingeckoPriceFeed(tx, coingeckoPriceFeed); err != nil {
			return err
		}
	} else {
		existingCoingeckoPriceFeed.Price = coingeckoPriceFeed.Price
		existingCoingeckoPriceFeed.Currency = coingeckoPriceFeed.Currency
		existingCoingeckoPriceFeed.UpdatedAt = time.Now()
		if err := UpdateCoingeckoPriceFeed(tx, existingCoingeckoPriceFeed); err != nil {
			return err
		}
	}

	return nil
}

// GetCoingeckoPriceFeedByName returns a CoingeckoPriceFeed by its name
func GetCoingeckoPriceFeedByName(tx *gorm.DB, name string) (*CoingeckoPriceFeed, error) {
	coingeckoPriceFeed := CoingeckoPriceFeed{}
	err := tx.Where("name = ?", name).First(&coingeckoPriceFeed).Error

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

// UpdateCoingeckoPriceFeed updates a CoingeckoPriceFeed
func UpdateCoingeckoPriceFeed(tx *gorm.DB, coingeckoPriceFeed *CoingeckoPriceFeed) error {
	if result := tx.Save(coingeckoPriceFeed); result.Error != nil {
		return result.Error
	}
	return nil
}
