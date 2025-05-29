package main

import (
	"os"

	"github.com/JonasLindermayr/FileBeam/handlers"
	"github.com/JonasLindermayr/FileBeam/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const Filepath = "data.db"
const UploadFilePath = "./uploads"

func main() {

	if _, err := os.Stat(Filepath); os.IsNotExist(err) {
		internal.Log("Database file not found. Please run migrate.go first.", internal.ERROR)
		panic(err)
	}
	if _, err := os.Stat(UploadFilePath); os.IsNotExist(err) {
		os.Mkdir(UploadFilePath, os.ModePerm)
	}

	var err error

	internal.DB, err = gorm.Open(sqlite.Open(Filepath), &gorm.Config{})
	if err != nil {
		internal.Log("Failed to connect to database", internal.ERROR)
		internal.Log("Error: "+err.Error(), internal.ERROR)
		panic(err)
	}
	internal.Log("Connected to database", internal.INFO)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	//gin.SetMode(gin.ReleaseMode)

	// OPEN ROUTES
	router.POST("/api/user", handlers.GetUserHandler)
	router.POST("/api/user/verify-otp", handlers.VerifyUserOTP)

	// CLOSED ROUTES -> MIDDLEWARE
	router.POST("/api/user/logout", internal.JWTAuthMiddleware(), handlers.LogoutUserHandler)
	router.POST("/api/upload", internal.JWTAuthMiddleware(), handlers.UploadHandler)
	router.Run(":8072")
}
