package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserId           uint      `gorm:"primaryKey"` // GORM automatically uses fields with name 'ID' as the primary key
	Username         string    `gorm:"size:255;not null"`
	Email            string    `gorm:"size:255;unique;not null"`
	HashedPublicKey  string    `gorm:"type:text"`
	OtherUserDetails string    `gorm:"type:text"`
	SignupDate       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func GetAllUsers(db *gorm.DB) []User {

	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		log.Fatalf("Error retrieving users: %v", result.Error)
	}

	// Print out the retrieved users
	// for _, user := range users {
	// 	fmt.Printf("User ID: %v,  User Name: %v, Email: %v\n", user.UserId, user.UserName, user.Email)
	// }
	return users
}
