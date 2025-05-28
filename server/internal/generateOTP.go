package internal

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)


func GenerateOTP() (string, error) {
	const otpLenght = 6
	var otp string

	for i := 0; i < otpLenght; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += num.String()
	}

	return otp, nil
}


func GenerateOTPRequestID() (string, error) {
	const requestIDLength = 16
	bytes := make([]byte, requestIDLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	otpRequestID := hex.EncodeToString(bytes)
	return otpRequestID, nil
}