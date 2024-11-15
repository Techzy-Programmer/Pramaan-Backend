package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"pramaan-chain/internal/db"
	"pramaan-chain/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ToDo: Add support for dynamic blob object names
var containerName = os.Getenv("BLOB_CONTAINER")

func UploadEvidenceHandler(c *gin.Context) {
	ext, ef := c.Request.Header["X-Evidence-Extension"]
	pubAddr, pf := c.Request.Header["X-Pub-Address"]
	hash, hf := c.Request.Header["X-Evidence-Hash"]

	if (!ef || !pf || !hf) || (len(hash[0]) != 128 || len(ext[0]) == 0) || ext[0][0] != '.' {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid evidence hash or extension",
		})
		return
	}

	success := false
	blobUploadPath := pubAddr[0] + "/" + hash[0]
	cErr := db.CreateInitialEvidenceRecord(pubAddr[0], hash[0], ext[0], blobUploadPath)
	if cErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Looks like this evidence already exists",
		})
		return
	}
	defer func() {
		if !success {
			db.DeleteEvidenceRecord(hash[0])
		}
	}()

	client, err := getBlobServiceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := c.FormFile("evidence")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve evidence",
		})
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open evidence stream",
		})
		return
	}
	defer srcFile.Close()

	_, upErr := client.UploadStream(context.TODO(), containerName, blobUploadPath, srcFile, nil)

	if upErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload evidence to Secure Storage",
		})
		return
	}

	success = true
	c.JSON(http.StatusOK, gin.H{
		"message": "Evidence uploaded successfully",
	})
}

func ConfirmedEvidenceHandler(c *gin.Context) {
	ctx, cf := c.Request.Header["X-Evidence-Creation-Tx"]
	pubAddr, pf := c.Request.Header["X-Pub-Address"]
	hash, hf := c.Request.Header["X-Evidence-Hash"]

	if !pf || !hf || !cf || len(hash[0]) != 128 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid evidence hash or public address",
		})
		return
	}

	evidence, dbErr := db.RetrieveEvidenceRecord(pubAddr[0], hash[0])
	if dbErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Evidence with provided hash doesn't exists",
		})
		return
	}

	index, iErr := strconv.Atoi(c.Param("index"))
	if iErr != nil || index < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid evidence index",
		})
		return
	}

	evidence.Index = index
	evidence.CreationTx = ctx[0]
	uErr := db.UpdateEvidenceRecord(evidence)
	if uErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update evidence record",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Evidence confirmed successfully",
	})
}

func DownloadEvidenceHandler(c *gin.Context) {
	accessSig, af := c.Request.Header["X-Access-Signature"]
	selfPubAddr, pf := c.Request.Header["X-Pub-Address"]
	evHash := c.Param("evHash") // SHA-512 Hash
	masterPubAddr := c.Param("pubAddr")

	if !pf || !af || len(evHash) != 128 || len(masterPubAddr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid evidence hash or public address or access signature",
		})
		return
	}

	evidence, rErr := db.RetrieveEvidenceRecord(masterPubAddr, evHash)
	if rErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Evidence with provided hash doesn't exists",
		})
		return
	}

	if evidence.OwnerAddr != selfPubAddr[0] {
		subOwner, dbErr := db.RetrieveOwner(selfPubAddr[0])
		if dbErr != nil || *subOwner.MasterId != masterPubAddr {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You don't have access to this evidence",
			})
			return
		}

		verified := utils.VerifySignature(masterPubAddr, strconv.Itoa(int(*subOwner.AccessTimestamp)), accessSig[0])
		if !verified {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Signature verification failed",
			})
			return
		}

		if int64(*subOwner.AccessTimestamp) < time.Now().Unix() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Your access for this evidence has expired",
			})
			return
		}
	}

	client, err := getBlobServiceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	blobDownloadPath := masterPubAddr + "/" + evHash
	downloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobDownloadPath, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download evidence from Blob Storage",
		})
		return
	}
	defer downloadResponse.Body.Close()

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+evHash+evidence.Extension)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Writer, downloadResponse.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to stream evidence for download",
		})
		return
	}
}

func ListEvidencesHandler(c *gin.Context) {
}
