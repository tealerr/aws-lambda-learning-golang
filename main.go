package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	ddlambda "github.com/DataDog/dd-trace-go/contrib/aws/datadog-lambda-go/v2"
	"github.com/aws/aws-lambda-go/events"
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

func newHandler(router *gin.Engine) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	adapter := ginadapter.New(router)
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return adapter.ProxyWithContext(ctx, event)
	}
}

func main() {
	router := setupRouter()
	// Run as a local HTTP server when no Lambda Runtime API is available.
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
	handler := newHandler(router)
	lambda.Start(ddlambda.WrapFunction(handler, nil))
}
