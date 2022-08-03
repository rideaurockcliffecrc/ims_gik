package middleware

import (
	"GIK_Web/database"
	"GIK_Web/types"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("userId").(uint)

		// get the user
		user := types.User{}
		if err := database.Database.Where("id = ?", userId).First(&user).Error; err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Invalid user",
			})
			c.Abort()
			return
		}

		if !user.Admin {
			c.JSON(401, gin.H{
				"success": false,
				"message": "You are not permitted to do this",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
