package handler

import (
	"fmt"
	"net/http"
	"os"
	"pramaan-chain/internal/db"
	"pramaan-chain/utils"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

func getBlobServiceClient() (*azblob.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
	}

	client, err := azblob.NewClient(os.Getenv("BLOB_URL"), cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %v", err)
	}

	return client, nil
}

func getPubAddress(c *gin.Context) (string, error) {
	pubAddr, addrFound := c.Request.Header["X-Pub-Address"]
	if !addrFound {
		return "", fmt.Errorf("X-Pub-Address header not found")
	}

	return pubAddr[0], nil
}

func verifyEvidenceAccess(c *gin.Context, selfPubAddr string, masterPubAddr string) bool {
	accessSig, af := c.Request.Header["X-Access-Signature"]
	if !af {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access signature not provided",
		})
		return false
	}

	subOwner, dbErr := db.RetrieveOwner(selfPubAddr)
	if dbErr != nil || *subOwner.MasterId != masterPubAddr {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this evidence",
		})
		return false
	}

	verified := utils.VerifySignature(masterPubAddr, strconv.Itoa(int(*subOwner.AccessTimestamp)), accessSig[0])
	if !verified {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Signature verification failed",
		})
		return false
	}

	if int64(*subOwner.AccessTimestamp) < time.Now().Unix() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Your access for this evidence has expired",
		})
		return false
	}

	return true
}
