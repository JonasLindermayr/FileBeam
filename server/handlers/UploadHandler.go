package handlers

import (
	"net/http"

	"github.com/JonasLindermayr/FileBeam/handlers/controller"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadHandler(c *gin.Context) {
	uuidValue, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userUUID := uuidValue.(uuid.UUID)

	fileID, fileSize, ext, err := controller.UploadFile(c, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"fileId":   fileID,
		"fileSize": fileSize,
		"ext":      ext,
	})
}
