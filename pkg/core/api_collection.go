/*
 * Paged
 *
 * Handles CRUD operations for events
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddEvent - 
func AddEvent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// GetEvents - 
func GetEvents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
