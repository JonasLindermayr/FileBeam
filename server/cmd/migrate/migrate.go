package main

import (
	"os"

	"github.com/JonasLindermayr/FileBeam/handlers/controller"
	"github.com/JonasLindermayr/FileBeam/internal"
	"github.com/JonasLindermayr/FileBeam/types"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	var err error
	internal.DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	internal.LogMigrate("Database initialized", internal.INFO)

	err = internal.DB.AutoMigrate(&types.User{}, &types.Session{})
	if err != nil {
		panic(err)
	}
	internal.LogMigrate("Database migrated", internal.WARNING)

	controller.CreateUserWithMigrate("admin", "password", "6r6E6@example.com")
	internal.LogMigrate("Migration completed", internal.INFO)
	os.Exit(0)
}
