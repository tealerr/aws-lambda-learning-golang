package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	return r
}

func main() {
	router := setupRouter()
	// Run as a local HTTP server when no Lambda Runtime API is available. .
	// LOCAL_SERVER=true remains an explicit override for local development.
	if strings.EqualFold(os.Getenv("LOCAL_SERVER"), "true") || os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		const address = ":8080"
		log.Printf("🚀 Local server running at http://localhost%s", address)
		log.Printf("📍 Endpoints: GET / (Hello World), GET /health (Health check)")
		if err := router.Run(address); err != nil {
			panic(err)
		}
		return
	}
	adapter := ginadapter.New(router)
	lambda.Start(adapter.Proxy)
}
