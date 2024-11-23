package main

import (
	"net/http"

	"github.com/adedejiosvaldo/scalable_api/db"
	"github.com/adedejiosvaldo/scalable_api/handlers"
	"github.com/gin-gonic/gin"
)

var requestCount uint64 // A counter to track the number of requests

func main() {
	db.ConnectPostgres()
	db.ConnectMongo()
	// db.ConnectRedis()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "Running",
		})
	})

	r.GET("/benchmark/sql/write", handlers.SQLWriteBenchmarkHandler)
	r.GET("/benchmark/sql/read", handlers.SQLReadBenchmarkHandler)
	// Add MongoDB and Redis handlers here

	r.Run(":8080")
}

// func main() {
// 	r := gin.Default()

// 	// Route to handle requests
// 	r.GET("/request", func(c *gin.Context) {
// 		// Increment the counter atomically
// 		currentRequest := atomic.AddUint64(&requestCount, 1)

// 		// Generate response
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": fmt.Sprintf("You are request number %d", currentRequest),
// 		})
// 	})

// 	// Start the server on port 8080
// 	r.Run(":8080")
// }
