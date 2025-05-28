package controller

import (
	"net/http"
	"time"

	"github.com/JonasLindermayr/FileBeam/internal"
	"github.com/JonasLindermayr/FileBeam/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *gin.Context) {

	var authInput types.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		internal.Log("Failed to bind JSON", internal.ERROR)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var userFound types.User
	internal.DB.Where("username = ?", authInput.Username).First(&userFound)

	if userFound.ID == uuid.Nil {
		internal.Log("User not found", internal.ERROR)
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		internal.Log("Invalid password", internal.ERROR)
		internal.Log("Username: "+authInput.Username, internal.WARNING)
		internal.Log("IP: "+c.ClientIP(), internal.DEBUG)
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid password"})
		return
	}

	otp, err := internal.GenerateOTP()
	if err != nil {
		internal.Log("Failed to generate OTP", internal.ERROR)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate OTP"})
		return
	}

	otpRequestId, err := internal.GenerateOTPRequestID()
	if err != nil {
		internal.Log("Failed to generate OTP request ID", internal.ERROR)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate OTP request ID"})
		return
	}

	userFound.OTP = otp
	userFound.OTPRequestID = otpRequestId
	internal.DB.Save(&userFound)
	internal.Log("Login request for uuid: "+userFound.ID.String(), internal.INFO)
	internal.Log("Created new OTP: "+userFound.OTP, internal.INFO)

	c.JSON(http.StatusOK, gin.H{
		"userId":       userFound.ID,
		"otpRequestID": otpRequestId,
	})
}

func VerifyOTP(c *gin.Context) {

	var authInput types.OTPAuthInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		internal.Log("Failed to bind JSON", internal.ERROR)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var userFound types.User
	internal.DB.Where("ID = ?", authInput.UUID).First(&userFound)

	if userFound.ID == uuid.Nil {
		internal.Log("User not found", internal.ERROR)
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		return
	}

	if userFound.OTPRequestID != authInput.OTPRequestID {
		internal.Log("Invalid OTP request ID", internal.ERROR)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid OTP request ID"})
		return
	}

	if userFound.OTP != authInput.OTP {
		internal.Log("Invalid OTP", internal.ERROR)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid OTP"})
		return
	}

	token, err := internal.GenerateToken(userFound.ID, time.Now().Add(time.Hour*24*30))
	if err != nil {
		internal.Log("Failed to generate token", internal.ERROR)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate token"})
		return
	}
	internal.Log("Generated token", internal.INFO)
	var session = &types.Session{
		Token:      token,
		EmployeeID: userFound.ID,
	}
	internal.DB.Create(&session)
	internal.Log("Session created for uuid: "+userFound.ID.String(), internal.INFO)

	internal.Log("User logged in as: "+userFound.Username, internal.INFO)
	userFound.OTP = ""
	userFound.OTPRequestID = ""
	internal.DB.Save(&userFound)
	internal.Log("OTP + OTPRequestID cleared", internal.INFO)

	c.JSON(http.StatusOK, gin.H{
		"userId":   userFound.ID,
		"username": userFound.Username,
		"email":    userFound.Email,
		"token":    token,
	})

}

func CreateUser(username, password, email string) {

	authInput := &types.AuthInput{
		Username: username,
		Password: password,
	}

	var userFound types.User
	internal.DB.Where("username = ?", authInput.Username).First(&userFound)

	if userFound.ID != uuid.Nil {
		internal.Log("User already exists", internal.ERROR)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		internal.Log("Failed to hash password", internal.ERROR)
	}

	user := &types.User{
		ID:       uuid.New(),
		Username: authInput.Username,
		Password: string(passwordHash),
		Email:    email,
	}

	internal.DB.Create(&user)
	internal.Log("User created", internal.INFO)

}

func CreateUserWithMigrate(username, password, email string) {

	authInput := &types.AuthInput{
		Username: username,
		Password: password,
	}

	var userFound types.User
	internal.DB.Where("username = ?", authInput.Username).First(&userFound)

	if userFound.ID != uuid.Nil {
		internal.LogMigrate("User already exists", internal.ERROR)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		internal.LogMigrate("Failed to hash password", internal.ERROR)
		return
	}

	user := &types.User{
		ID:       uuid.New(),
		Username: authInput.Username,
		Password: string(passwordHash),
		Email:    email,
	}

	internal.DB.Create(&user)
	internal.LogMigrate("User created", internal.INFO)

}

func LogoutUser(c *gin.Context) {
	uuidValue, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userUUID := uuidValue.(string)

	var sessions []types.Session
	internal.DB.Where("EmployeeID = ?", userUUID).Find(&sessions)

	for _, session := range sessions {
		internal.DB.Delete(&session)
	}

	internal.Log("User logged out", internal.INFO)
	c.JSON(http.StatusOK, gin.H{"message": "User logged out"})
}
