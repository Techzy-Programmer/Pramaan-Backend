package handler

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

func VerifySignature(c *gin.Context) {
	pubAddr, addrFound := c.Request.Header["X-Pub-Address"]
	if !addrFound {
		c.JSON(403, gin.H{"error": "X-Pub-Address header not found"})
		c.Abort()
		return
	}

	sig, sigFound := c.Request.Header["X-Signature"]
	if !sigFound {
		c.JSON(403, gin.H{"error": "X-Signature header not found"})
		c.Abort()
		return
	}

	if !verifySignature(pubAddr[0], sig[0]) {
		c.JSON(403, gin.H{"error": "Invalid signature"})
		c.Abort()
		return
	}

	c.Next()
}

func verifySignature(pubAddr string, sig string) bool {
	message := "Authorize Me!"
	hash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	sigBytes := common.FromHex(sig)
	if len(sigBytes) != 65 {
		return false
	}

	r := sigBytes[:32]
	s := sigBytes[32:64]
	v := sigBytes[64]

	if v == 27 || v == 28 {
		v -= 27
	}

	publicKey, err := crypto.SigToPub(hash.Bytes(), append(r, append(s, v)...))
	if err != nil {
		fmt.Printf("Failed to recover public key: %v\n", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*publicKey)
	return bytes.Equal(recoveredAddress.Bytes(), common.HexToAddress(pubAddr).Bytes())
}

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
