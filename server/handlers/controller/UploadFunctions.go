package controller

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JonasLindermayr/FileBeam/internal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context, userUUID uuid.UUID) (string, int64, string, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return "", 0, "", fmt.Errorf("no file found")
	}
	defer file.Close()

	userFolder := fmt.Sprintf("./uploads/%s", userUUID)

	if err := os.MkdirAll(userFolder, os.ModePerm); err != nil {
		return "", 0, "", fmt.Errorf("could not create user directory")
	}

	ext := filepath.Ext(header.Filename)
	fileID := uuid.New().String()
	dst := fmt.Sprintf("%s/%s%s", userFolder, fileID, ext)

	if err := c.SaveUploadedFile(header, dst); err != nil {
		return "", 0, "", fmt.Errorf("could not save file")
	}

	internal.Log(fmt.Sprintf("User -> %s: uploaded File -> %s with a size of: %s",
		userUUID.String(), fileID, internal.FormatBytes(header.Size)), internal.INFO)

	return fileID, header.Size, ext, nil
}
