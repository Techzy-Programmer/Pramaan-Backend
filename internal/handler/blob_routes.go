package handler

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ToDo: Add support for dynamic blob object names
var containerName = os.Getenv("BLOB_CONTAINER")

func UploadHandler(c *gin.Context) {
	client, err := getBlobServiceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve file",
		})
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file",
		})
		return
	}
	defer srcFile.Close()

	blobName := "uploaded-file.mp4"
	_, upErr := client.UploadStream(context.TODO(), containerName, blobName, srcFile, nil)

	if upErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file to Blob Storage > " + upErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}

func DownloadHandler(c *gin.Context) {
	client, err := getBlobServiceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	blobName := "uploaded-file.mp4"
	downloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobName, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file from Blob Storage"})
		return
	}
	defer downloadResponse.Body.Close()

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+blobName)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Writer, downloadResponse.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file to response"})
		return
	}
}
