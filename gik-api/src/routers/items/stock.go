package items

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"GIK_Web/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func JumpStock(c *gin.Context) {
	difference := c.Query("diff")
	productId := c.Query("product_id")

	// parse difference
	differenceInt, err := strconv.Atoi(difference)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid difference",
		})
		return
	}

	// get item
	item := types.Item0{}
	err = database.Database.Model(&types.Item0{}).Where("product_id = ?", productId).Scan(&item).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid product",
		})
		return
	}

	err = database.Database.Model(&types.Item0{}).Where("product_id = ?", productId).Update("stock_quantity", item.StockQuantity+float32(differenceInt)).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to update stock",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Stock updated",
		"data":    item.StockQuantity + float32(differenceInt),
	})

	utils.CreateSimpleLog(c, fmt.Sprintf("Jumped stock for product ID: %s by %d to %f", productId, differenceInt, item.StockQuantity+float32(differenceInt)))
}

func AddStock(c *gin.Context) {
	productId := c.Query("product_id")

	if productId == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	item := types.Item0{}
	err := database.Database.Model(&types.Item0{}).Where("product_id = ?", productId).Scan(&item).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid product",
		})
		return
	}

	err = database.Database.Model(&types.Item0{}).Where("product_id = ?", productId).Update("stock_quantity", item.StockQuantity+1).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to update stock",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Stock +1 updated",
	})

	utils.CreateSimpleLog(c, "Added stock for product ID: "+productId)

}

func RemoveStock(c *gin.Context) {
	productId := c.Query("product_id")

	if productId == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	item := types.Item0{}
	err := database.Database.Model(&types.Item0{}).Where("product_id = ?", productId).Scan(&item).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid product",
		})
		return
	}

	// update
	err = database.Database.Model(&types.Item0{}).Where("product_id = ?", productId).Update("stock_quantity", item.StockQuantity-1).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to update stock",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Stock -1 updated",
	})

	utils.CreateSimpleLog(c, "Removed stock for product ID: "+productId)

}
