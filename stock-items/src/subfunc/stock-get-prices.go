package subfunc

import "github.com/gin-gonic/gin"

// GetHistoricalPrices ...
func GetHistoricalPrices(c *gin.Context) {
	text := "test2"
	c.JSON(200, gin.H{
		"body": text,
	})
}
