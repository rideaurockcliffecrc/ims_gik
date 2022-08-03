package admin

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ToggleSignupCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid sign up code",
		})
		return
	}

	// try to find signup code
	signupCode := types.SignupCode{}
	if err := database.Database.Where(&types.SignupCode{
		Code: code,
	}).First(&signupCode).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid sign up code",
		})
		return
	}

	if signupCode.Expiration < time.Now().Unix() && !signupCode.Expired {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Sign up code has expired",
		})
		return
	}

	signupCode.Expired = !signupCode.Expired

	if err := database.Database.Save(&signupCode).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Error saving sign up code",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Sign up code toggled",
		"data":    signupCode.Expired,
	})
}

func CreateSignupCode(c *gin.Context) {

	designatedUsername := c.Query("username")

	if designatedUsername == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid username designation",
		})
		return
	}

	// make sure no other account has this username
	var count int64
	database.Database.Model(&types.User{}).Where("username = ?", designatedUsername).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Username already exists",
		})
		return
	}

	var count2 int64
	database.Database.Model(&types.SignupCode{}).Where(&types.SignupCode{
		DesignatedUsername: designatedUsername,
	}).Count(&count2)

	if count2 > 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Signup code already exists",
		})
		return
	}

	newSignupCode := types.SignupCode{
		Code:               uuid.New().String(),
		Expiration:         time.Now().Unix() + (60 * 60 * 24 * 7),
		DesignatedUsername: designatedUsername,
		//CreatedByUserID:    c.MustGet("userId").(uint),
		Expired:   false,
		CreatedAt: time.Now().Unix(),
	}

	if err := database.Database.Create(&newSignupCode).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to create signup code",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Signup code created",
		"data":    newSignupCode.Code,
	})
}

func GetSignupCodes(c *gin.Context) {
	// paginate signup codes

	page := c.Query("page")

	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid page number",
		})
		return
	}

	var signupCodes []types.SignupCode

	limit := 10
	offset := (pageInt - 1) * limit

	baseQuery := database.Database.Model(&types.SignupCode{})

	var totalCount int64
	baseQuery.Count(&totalCount)

	baseQuery = baseQuery.Limit(limit).Offset(offset).Order("created_at desc")

	baseQuery.Find(&signupCodes)

	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"data":        signupCodes,
			"total":       totalCount,
			"currentPage": pageInt,
			"totalPages":  totalPages,
		},
	})
}
