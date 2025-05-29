package types

import "github.com/google/uuid"

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPAuthInput struct {
	UUID         string `json:"uuid" binding:"required"`
	OTP          string `json:"otp" binding:"required"`
	OTPRequestID string `json:"otpRequestID" binding:"required"`
}

type JWT struct {
	UUID uuid.UUID `json:"uuid" binding:"required"`
}
