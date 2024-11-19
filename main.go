package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var requestCount uint64 // A counter to track the number of requests

func main() {
	r := gin.Default()

	// Route to handle requests
	r.GET("/request", func(c *gin.Context) {
		// Increment the counter atomically
		currentRequest := atomic.AddUint64(&requestCount, 1)

		// Generate response
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("You are request number %d", currentRequest),
		})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
