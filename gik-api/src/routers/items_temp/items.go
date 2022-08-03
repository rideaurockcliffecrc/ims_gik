package items_temp

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"GIK_Web/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type item struct {
	ID       string
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Category string  `json:"category"` //Gender and such
	Size     string  `json:"size"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}

type returnedItem struct {
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Category string  `json:"category"` //Gender and such
	Size     string  `json:"size"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
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

	baseQuery := database.Database.Model(&types.Item2{})

	baseQuery = baseQuery.Order("sku, FIELD(size, 'XXL',  'XL', 'L', 'M', 'S', 'XS', 'XXS'), size")

	if search != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+search+"%")
	}

	var totalCount int64
	baseQuery.Count(&totalCount)

	baseQuery = baseQuery.Limit(limit).Offset(offset)

	items := []item{}

	baseQuery.Find(&items)

	//returnedItems := []returnedItem{}
	/*
		for _, item := range items {
			itemName := item.Name

				if item.Name == "" {
					tempItemName, err := utils.GetItemNameByID(item.ID)
					if tempItemName == "" || err != nil {
						fmt.Println("unable to find name for", item.SKU)
						tempItemName = "Unknown name"
					}

					itemName = tempItemName

					go database.Database.Model(&types.Item{}).Where("product_id = ?", item.ProductID).Update("name", itemName)
				}

			returnedItems = append(returnedItems, returnedItem{
				Name:          itemName,
				SKU:           item.SKU,
			})
		}*/

	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	c.JSON(200, gin.H{"success": true, "data": gin.H{
		"data":        items,
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

	item := types.Item2{}
	if err := database.Database.Where(types.Item{
		ProductID: uint(jsonIdInt),
	}).First(&item).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	item.Name = json.Name
	item.SKU = json.SKU
	item.Category = json.Category
	item.Size = json.Size
	item.Price = json.Price
	item.Quantity = json.Quantity

	database.Database.Save(item)
}

type newItemRequest struct {
	Name     string  `json:"name" binding:"required"`
	SKU      string  `json:"sku" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Size     string  `json:"size" binding:"required"`
	Price    float32 `json:"price" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
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

	item := types.Item2{}

	item.Name = json.Name
	item.SKU = json.SKU
	item.Category = json.Category
	item.Size = json.Size
	item.Price = json.Price
	item.Quantity = int(json.Quantity)

	err := database.Database.Create(&item).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to create item",
			"error":   err.Error(),
		})
		return
	}

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
