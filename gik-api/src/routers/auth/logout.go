package auth

import (
	"GIK_Web/env"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", env.CookieDomain, false, false)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}
