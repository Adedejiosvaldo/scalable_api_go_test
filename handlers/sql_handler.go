package handlers

import (
	"net/http"
	"strconv"

	"github.com/adedejiosvaldo/scalable_api/db"
	"github.com/gin-gonic/gin"
)

func WriteSQLHandler(c *gin.Context) {
	data := c.Query("data")

	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data is required"})
		return
	}
	err := db.InsertToPostgres(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}

func ReadSQLHandler(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	data, err := db.ReadFromPostgres(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
