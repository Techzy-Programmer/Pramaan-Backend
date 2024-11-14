package server

import (
	"fmt"
	"net/http"
	"pramaan-chain/internal/handler"

	"github.com/gin-gonic/gin"
)

func StartAPIServer(port string) {
	router := gin.Default()
	router.Use(handler.VerifySignature) // Web3 Signature Verification Middleware

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
	// Blob Routes
	route.POST("/upload", handler.UploadHandler)
	route.GET("/download", handler.DownloadHandler)
}
