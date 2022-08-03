package items_temp

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAutoSuggest(c *gin.Context) {
	query := c.Query("query")

	items := []types.Item2{}
	err := database.Database.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query)).Find(&items).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to query items",
		})
		return
	}

	suggestions := []string{}
	for _, item := range items {
		suggestions = append(suggestions, item.Name)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Auto suggested items",
		"data":    suggestions,
	})
}
