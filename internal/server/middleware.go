package server

import (
	"pramaan-chain/utils"

	"github.com/gin-gonic/gin"
)

func validateSignature(c *gin.Context) {
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

	if !utils.VerifySignature(pubAddr[0], "Authorize Me!", sig[0]) {
		c.JSON(403, gin.H{"error": "Invalid signature"})
		c.Abort()
		return
	}

	c.Next()
}
