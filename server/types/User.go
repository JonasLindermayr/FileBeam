package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Username     string    `json:"username" gorm:"unique"`
	Password     string    `json:"password"`
	Email        string    `json:"email" gorm:"unique"`
	OTP          string    `json:"otp" gorm:"unique"`
	OTPRequestID string    `json:"otpRequestID"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
