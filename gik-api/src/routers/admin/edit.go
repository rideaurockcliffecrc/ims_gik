package admin

import (
	"GIK_Web/database"
	"GIK_Web/types"

	"github.com/gin-gonic/gin"
)

func EditAdmins(c *gin.Context) {
	json := [][]formattedUser{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Error binding JSON",
		})
		return
	}

	if len(json) != 2 {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Invalid JSON",
		})
		return
	}

	nonAdmins := json[0]
	admins := json[1]

	for _, admin := range admins {
		// set these users to admin
		if err := database.Database.Model(&types.User{}).Where("id = ?", admin.ID).Update("admin", true).Error; err != nil {
			continue
		}
	}

	for _, nonAdmin := range nonAdmins {
		// set these users to non-admin
		if err := database.Database.Model(&types.User{}).Where("id = ?", nonAdmin.ID).Update("admin", false).Error; err != nil {
			continue
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Admin list updated",
	})

}
