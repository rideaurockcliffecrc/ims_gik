package middleware

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"time"

	"github.com/gin-gonic/gin"
)

func AdvancedLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := c.MustGet("userId").(uint)

		if user == 0 {
			c.Next()
			return
		}

		newLog := types.AdvancedLog{
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			UserID:    user,
			Timestamp: time.Now().Unix(),
		}

		if err := database.Database.Create(&newLog).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Unable to log action",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
