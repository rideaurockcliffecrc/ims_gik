package items

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"GIK_Web/utils"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type item struct {
	ID            string
	SKU           string  `json:"sku"`
	ProductID     uint    `json:"id"`
	Price         float32 `json:"price"`
	StockQuantity float32 `json:"stockQuantity"`
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Gender        string  `json:"gender"`
	Season        string  `json:"season"`
	Location      string  `json:"location"`
	Image         []byte  `json:"image"`
}

type returnedItem struct {
	ProductID     uint    `json:"id"`
	StockQuantity float32 `json:"stockQuantity"`
	Name          string  `json:"name"`
	SKU           string  `json:"sku"`
	Gender        string  `json:"gender"`
	Season        string  `json:"season"`
	Price         float32 `json:"price"`
}

func ListItem(c *gin.Context) {
	page := c.Query("page")

	if page == "" {
		page = "1"
	}

	search := c.Query("search")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid page number",
		})
		return
	}

	limit := 10
	offset := (pageInt - 1) * limit

	baseQuery := database.Database.Model(&types.Item{})

	baseQuery = baseQuery.Order("timestamp desc")

	if search != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+search+"%")
	}

	var totalCount int64
	baseQuery.Count(&totalCount)

	baseQuery = baseQuery.Limit(limit).Offset(offset)

	items := []item{}

	baseQuery.Find(&items)

	returnedItems := []returnedItem{}

	for _, item := range items {
		itemName := item.Name

		if item.Name == "" {
			tempItemName, err := utils.GetItemNameByID(item.ProductID)
			if tempItemName == "" || err != nil {
				fmt.Println("unable to find name for", item.ProductID)
				tempItemName = "Unknown name"
			}

			itemName = tempItemName

			go database.Database.Model(&types.Item{}).Where("product_id = ?", item.ProductID).Update("name", itemName)
		}

		returnedItems = append(returnedItems, returnedItem{
			ProductID:     item.ProductID,
			StockQuantity: item.StockQuantity,
			Name:          itemName,
			SKU:           item.SKU,
		})
	}

	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	c.JSON(200, gin.H{"success": true, "data": gin.H{
		"data":        returnedItems,
		"total":       totalCount,
		"currentPage": pageInt,
		"totalPages":  totalPages,
	}})

}

func UpdateItem(c *gin.Context) {
	json := item{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	jsonIdInt, err := strconv.Atoi(json.ID)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	item := types.Item{}
	if err := database.Database.Where(types.Item{
		ProductID: uint(jsonIdInt),
	}).First(&item).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	item.Price = json.Price
	item.StockQuantity = json.StockQuantity
	item.Name = json.Name
	item.Category = json.Category
	item.Gender = json.Gender
	item.Season = json.Season
	item.Location = json.Location
	item.Image = json.Image

	database.Database.Save(item)
}

type newItemRequest struct {
	Name          string  `json:"name" binding:"required"`
	SKU           string  `json:"sku" binding:"required"`
	StockQuantity int     `json:"stockQuantity" binding:"required"`
	Gender        string  `json:"gender" binding:"required"`
	Season        string  `json:"season" binding:"required"`
	Price         float32 `json:"price" binding:"required"`
}

func AddItem(c *gin.Context) {
	json := newItemRequest{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	item := types.Item{}

	item.Name = json.Name
	item.SKU = json.SKU
	item.StockQuantity = float32(json.StockQuantity)
	item.ProductID = uint(utils.GenerateProductId())
	item.Timestamp = time.Now().Unix()
	item.Gender = json.Gender
	item.Season = json.Season
	item.Price = json.Price

	err := database.Database.Create(&item).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to create item",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"success": true, "data": gin.H{
		"id": item.ProductID,
	}})

	utils.CreateSimpleLog(c, fmt.Sprintf("Added item %s", item.Name))
}

func DeleteLocation(c *gin.Context) {
	json := item{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	database.Database.Where("ID = ?", json.ID).Where("name = ?", json.Name).Delete(&item{})
}
