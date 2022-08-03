package transaction

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ListTransactions(c *gin.Context) {
	transactions := []types.Transaction{}

	page := c.Query("page")

	limit := 10
	offset := 0

	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid page",
		})
		return
	}

	// pagination
	offset = (pageInt - 1) * limit

	baseQuery := database.Database.Model(&types.Transaction{}).Order("timestamp desc")

	totalCount := int64(0)
	baseQuery.Count(&totalCount)

	baseQuery.Limit(limit).Offset(offset).Find(&transactions)

	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"data":       transactions,
			"totalPages": totalPages,
		},
	})
}

type addRequest struct {
	ClientID int           `json:"clientId" binding:"required"`
	Type     bool          `json:"type"`
	Products []productBody `json:"products" binding:"required"`
}

type productBody struct {
	ID       int `json:"id" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

func AddTransaction(c *gin.Context) {
	json := addRequest{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	totalQuantity := 0

	for _, product := range json.Products {
		item := types.Item{}
		database.Database.Where("product_id = ?", product.ID).First(&item)

		if item.ProductID == 0 {
			c.JSON(400, gin.H{
				"success": false,
				"message": fmt.Sprintf("%d is an invalid product ID", product.ID),
			})
			return
		}

	}

	transaction := types.Transaction{
		Type:      json.Type,
		ClientID:  uint(json.ClientID),
		SignerID:  c.MustGet("userId").(uint),
		Timestamp: time.Now().Unix(),
	}

	database.Database.Create(&transaction)

	for _, product := range json.Products {
		// get items
		item := types.Item{}
		database.Database.Where("product_id = ?", product.ID).First(&item)

		// create transaction item
		transactionItem := types.TransactionItem{
			TransactionID: transaction.ID,
			Quantity:      product.Quantity,
			ProductID:     product.ID,
		}

		database.Database.Create(&transactionItem)

		totalQuantity += product.Quantity
	}

	transaction.TotalQuantity = totalQuantity

	database.Database.Save(&transaction)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Transaction created",
	})
}

func DeleteTransaction(c *gin.Context) {
	// get id
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	// get transaction
	transaction := types.Transaction{}
	database.Database.Where("id = ?", idInt).First(&transaction)

	if transaction.ID == 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Transaction not found",
		})
		return
	}

	// delete all transaction items
	transactionItems := []types.TransactionItem{}
	database.Database.Where("transaction_id = ?", transaction.ID).Delete(&transactionItems)

	// delete transaction
	database.Database.Delete(&transaction)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Transaction deleted",
	})
}

func GetTransactionItems(c *gin.Context) {
	// id
	id := c.Query("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	// get transaction
	transaction := types.Transaction{}
	database.Database.Where("id = ?", idInt).First(&transaction)

	if transaction.ID == 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Transaction not found",
		})
		return
	}

	// get transaction items
	transactionItems := []types.TransactionItem{}
	database.Database.Where("transaction_id = ?", transaction.ID).Find(&transactionItems)

	c.JSON(200, gin.H{
		"success": true,
		"data":    transactionItems,
	})
}
