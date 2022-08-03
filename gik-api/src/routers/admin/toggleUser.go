package admin

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ToggleUser(c *gin.Context) {
	userId := c.Query("user_id")

	if userId == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user id",
		})
		return
	}

	// convert to int
	userIdInt, err := strconv.Atoi(userId)
	if err != nil || userIdInt < 1 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user id",
		})
		return
	}

	// try to find user
	user := types.User{}
	if err := database.Database.Where("id = ?", userIdInt).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid user id",
		})
		return
	}

	user.Disabled = !user.Disabled

	if err := database.Database.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Error saving user",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "User toggled",
		"data":    user.Disabled,
	})

}
