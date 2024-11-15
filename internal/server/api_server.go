package server

import (
	"fmt"
	"net/http"
	"pramaan-chain/internal/handler"

	"github.com/gin-gonic/gin"
)

func StartAPIServer() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(validateSignature) // Web3 Signature Verification Middleware

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Wohoo! Nothing to be found here",
			"ok":      false,
		})
	})

	baseRoute := router.Group("/")
	registerEndpoints(*baseRoute)

	keyFile := "private.key"
	certFile := "fullchain.crt"
	fmt.Println("Serving backend on https port 443")

	if err := http.ListenAndServeTLS(":443", certFile, keyFile, router); err != nil {
		fmt.Println("> Stopping web server...\n", err.Error())
	}
}

func registerEndpoints(route gin.RouterGroup) {
	// Owner Routes
	route.POST("/grant-access", handler.GrantAccessHandler)
	route.POST("/owner/create", handler.CreateOwnerHandler)
	route.GET("/owner/retrieve", handler.RetrieveOwnerHandler)

	// Evidence Routes
	route.POST("/evidence/upload", handler.UploadEvidenceHandler)
	route.GET("/evidence/confirmed/:index", handler.ConfirmedEvidenceHandler)
	route.GET("/evidence/download/:evHash/:pubAddr", handler.DownloadEvidenceHandler)
}
