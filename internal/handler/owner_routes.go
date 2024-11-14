package handler

import (
	"pramaan-chain/internal/db"

	"github.com/gin-gonic/gin"
)

type CreateOwnerRequest struct {
	Name string `json:"name" binding:"required"`
	Tx   string `json:"tx" binding:"required"`
}

func CreateOwnerHandler(c *gin.Context) {
	var ownerReq CreateOwnerRequest

	bErr := c.ShouldBindJSON(&ownerReq)
	if bErr != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	pubAddr, pubErr := getPubAddress(c)
	if pubErr != nil {
		c.JSON(400, gin.H{
			"error": pubErr.Error(),
		})
		return
	}

	cErr := db.CreateOwner(&db.Owner{
		PubAddress: pubAddr,
		Name:       ownerReq.Name,
	})
	if cErr != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create owner record",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Owner record created successfully",
	})
}

func RetrieveOwnerHandler(c *gin.Context) {
	pubAddr, pubErr := getPubAddress(c)
	if pubErr != nil {
		c.JSON(400, gin.H{
			"error": pubErr.Error(),
		})
		return
	}

	owner, rErr := db.RetrieveOwner(pubAddr)
	if rErr != nil {
		c.JSON(404, gin.H{
			"error": "Owner record not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Owner record retrieved successfully",
		"owner":   owner,
	})
}

type GrantAccessRequest struct {
	SubOwnerPubAddr string `json:"subOwnerPubAddr" binding:"required"`
	MSG             string `json:"msg" binding:"required"`
	Tx              string `json:"tx" binding:"required"`
}

func GrantAccessHandler(c *gin.Context) {
	pubAddr, pubErr := getPubAddress(c)
	if pubErr != nil {
		c.JSON(400, gin.H{
			"error": pubErr.Error(),
		})
		return
	}

	var grantReq GrantAccessRequest
	bErr := c.ShouldBindJSON(&grantReq)
	if bErr != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	oErr := db.BridgeOwner(grantReq.SubOwnerPubAddr, grantReq.MSG, grantReq.Tx, pubAddr)
	if oErr != nil {
		c.JSON(500, gin.H{
			"error": "Failed to grant access",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Access granted successfully",
	})
}
