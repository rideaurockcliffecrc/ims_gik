package middleware

import (
	"GIK_Web/database"
	"GIK_Web/types"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the auth cookie
		authCookie, err := c.Cookie("session")

		// if the header is empty, return unauthorized
		if authCookie == "" || err != nil {
			c.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// attempt to find the session
		session := types.Session{}
		if err := database.Database.Where("id = ?", authCookie).First(&session).Error; err != nil {
			c.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		if err := database.Database.Where("id = ?", session.UserID).First(&types.User{}).Error; err != nil {
			c.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("userId", session.UserID)

		c.Next()
	}
}
