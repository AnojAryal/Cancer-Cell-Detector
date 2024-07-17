package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveImage(c *gin.Context, file *multipart.FileHeader, uploadDir string) (string, error) {

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate a new UUID
	uuidFileName := uuid.New().String() + filepath.Ext(file.Filename)

	filePath := filepath.Join(uploadDir, uuidFileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}
