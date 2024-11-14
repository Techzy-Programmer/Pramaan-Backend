package server

import (
	"fmt"
	"net/http"
	"pramaan-chain/internal/handler"

	"github.com/gin-gonic/gin"
)

func StartAPIServer(port string) {
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

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Println("Serving backend at :" + port)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Stopping web server...")
	}
}

func registerEndpoints(route gin.RouterGroup) {
	// Owner Routes
	route.POST("/grant-access", handler.GrantAccessHandler)
	route.POST("/owner/create", handler.CreateOwnerHandler)
	route.GET("/owner/retrieve", handler.RetrieveOwnerHandler)

	// Evidence Routes
	route.POST("/evidence/upload", handler.UploadEvidenceHandler)
	route.GET("/evidence/list/:pubAddr", handler.ListEvidencesHandler)
	route.GET("/evidence/download/:evId/:pubAddr", handler.DownloadEvidenceHandler)
}
