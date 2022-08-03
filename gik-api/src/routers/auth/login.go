package auth

import (
	"GIK_Web/database"
	"GIK_Web/env"
	"GIK_Web/types"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	//TotpCode        string `json:"totp"`
	VerificationJWT string `json:"verificationJWT" binding:"required"`
}

func Login(c *gin.Context) {
	json := loginRequest{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	// check if JWT is expired or invalid
	token, err := jwt.Parse(json.VerificationJWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.JWTSigningPassword), nil
	})
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid verification",
		})
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	// check if JWT is valid
	if !token.Valid {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid verification",
		})
		return
	}

	// check if JWT is expired
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Verification is expired",
		})
		return
	}

	// check if username matches up
	json.Username = strings.ToLower(json.Username)

	if claims["username"].(string) != json.Username {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Verification does not match",
		})
		return
	}

	// check if user exists
	user := types.User{}
	if err := database.Database.Model(&types.User{}).Where(&types.User{
		Username: json.Username,
	}).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid username or password",
			"code":    1,
		})
		return
	}
	/*
		// check if totp is correct
		if user.TwoFactorSecret != "" {
			// two factor is enabled, let's check if there's a code

			if json.TotpCode == "" {
				// must have a code
				c.JSON(400, gin.H{
					"success": false,
					"message": "This account has two factor enabled. Please provide a code.",
				})
				return
			}

			// check if code is valid

			if !totp.Validate(json.TotpCode, user.TwoFactorSecret) {
				c.JSON(400, gin.H{
					"success": false,
					"message": "Invalid two factor code",
				})
				return
			}
		}*/

	// create session
	session := types.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + (60 * 60 * 24),
	}

	if err := database.Database.Create(&session).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to create session",
		})
		return
	}

	c.SetCookie("session", session.ID, 60*60*24, "/", env.CookieDomain, false, true)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    session.ID,
	})

}
