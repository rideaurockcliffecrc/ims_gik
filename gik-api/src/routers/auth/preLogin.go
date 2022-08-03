package auth

import (
	"GIK_Web/database"
	"GIK_Web/env"
	"GIK_Web/types"
	"GIK_Web/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type checkRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CheckPasswordForLogin(c *gin.Context) {
	json := checkRequest{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	// check if user exists
	json.Username = strings.ToLower(json.Username)

	user := types.User{}
	if err := database.Database.Model(&types.User{}).Where(&types.User{
		Username: json.Username,
	}).First(&user).Error; err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid username or password",
			"code":    1,
		})
		return
	}

	if user.Disabled {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Account access is disabled. Contact an administrator to remove this restriction.",
		})
		return
	}

	// check if password is correct
	valid, err := utils.ComparePassword(json.Password, user.Password)
	if !valid || err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid username or password",
			"code":    2,
		})
		return
	}

	// now create a JWT to login (with 2fa)
	claims := &jwt.MapClaims{
		"iss":      "gikdash",
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
		"username": user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(env.JWTSigningPassword))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to create JWT",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Identity partially verified",
		"data":    tokenString,
	})

}
