package location

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"GIK_Web/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type location struct {
	ID     int
	Name   string `json:"name"`
	Letter string `json:"letter"`
	SKU    string `json:"sku"`
}

type lookupData struct {
	location
	Item types.Item `json:"product"`
}

type listData struct {
	location
	Item        types.Item `json:"product"`
	ProductName string     `json:"productName"`
}

func ListLocation(c *gin.Context) {
	name := c.Query("name")
	letter := c.Query("letter")
	productName := c.Query("productName")
	sku := c.Query("sku")

	locations := []location{}

	baseQuery := database.Database.Model(&types.Location{})

	if name != "" {
		baseQuery = baseQuery.Where("name = ?", name)
	}

	if letter != "" {
		baseQuery = baseQuery.Where("letter = ?", letter)
	}

	if productName != "" {

		skuAlt := ""

		database.Database.Model(&types.Item{}).Where("name = ?", productName).Distinct().Pluck("sku", &skuAlt)

		baseQuery = baseQuery.Where("sku = ?", skuAlt)
	}

	if sku != "" {
		baseQuery = baseQuery.Where("sku = ?", sku)
	}

	err := baseQuery.Find(&locations).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Unable to query locations",
		})
		return
	}

	response := []listData{}

	for _, location := range locations {
		var item types.Item
		err := database.Database.Model(&types.Item{}).Where("sku = ?", location.SKU).First(&item).Error

		if err != nil {
			continue
		}

		var name string

		database.Database.Model(&types.Item{}).Where("sku = ?", location.SKU).Distinct().Pluck("name", &name)

		response = append(response, listData{
			location:    location,
			Item:        item,
			ProductName: name,
		})

	}

	c.JSON(200, gin.H{"success": true, "data": response})
}

func UpdateLocation(c *gin.Context) {
	json := location{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	database.Database.Model(&location{}).Where("name = ?", json.Name).Update("sku", json.SKU)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Updated location",
	})

	utils.CreateSimpleLog(c, "Updated location "+json.Name)

}

type addRequest struct {
	Name   string `json:"name" binding:"required"`
	Letter string `json:"letter" binding:"required"`
	//SKU         string `json:"sku" binding:"required"`
	ProductName string `json:"productName" binding:"required"`
}

func AddLocation(c *gin.Context) {
	json := addRequest{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		fmt.Println(json)
		return
	}

	//TODO add check to see if exists

	sku := ""

	database.Database.Model(&types.Item{}).Where("name = ?", json.ProductName).Distinct().Pluck("sku", &sku)

	err := database.Database.Create(&types.Location{Name: json.Name, Letter: json.Letter, SKU: sku}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Unable to create location",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Location created",
	})

	utils.CreateSimpleLog(c, "Added location "+json.Name)

}

func DeleteLocation(c *gin.Context) {
	id := c.Query("id")

	// conver to integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	location := types.Location{}

	err = database.Database.Model(&location).Where("id = ?", idInt).Delete(&location).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Unable to delete location",
		})
		return
	}

	err = database.Database.Delete(&location).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Unable to delete location",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Location deleted",
	})

	utils.CreateSimpleLog(c, "Deleted location "+id)

}

func LookupLocation(c *gin.Context) {
	// product id
	name := c.Query("name")
	letter := c.Query("letter")

	//could remove check and use function also for list
	/*
		if name == "" && letter == "" && itemID == 0 {
			c.JSON(400, gin.H{
				"success": false,
				"message": "No fields provided",
			})
			return
		}*/

	var postData []location
	database.Database.Model(&location{}).Where(&location{Name: name, Letter: letter}).Scan(&postData)

	response := []lookupData{}

	for _, location := range postData {
		var item types.Item
		err := database.Database.Model(&types.Item{}).Where("sku = ?", location.SKU).Scan(&item).Error

		if err != nil {
			continue
		}

		response = append(response, lookupData{
			location: location,
			Item:     item,
		})

	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    response,
	})

}

func GetScannedData(c *gin.Context) {
	name := c.Query("name")
	letter := c.Query("letter")

	var product types.Item

	location := types.Location{}
	if err := database.Database.Model(&types.Location{}).Where(&types.Location{Name: name, Letter: letter}).Scan(&location).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Unable to query location",
			"data":    gin.H{},
		})
		return
	}

	var item types.Item
	if err := database.Database.Model(&types.Item{}).Where("sku = ?", location.SKU).First(&item).Error; err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Unable to find item",
			"error":   err.Error(),
			"data":    gin.H{},
		})
		return
	}

	product = item

	c.JSON(200, gin.H{
		"success": true,
		"data":    product,
	})

	utils.CreateSimpleLog(c, "Scanned location "+name+" "+letter)

}

func ListLocationSKU(c *gin.Context) {
	name := c.Query("name")
	letter := c.Query("letter")

	var sku string

	database.Database.Model(&types.Location{}).Where("name = ?", name).Where("letter = ?", letter).Distinct("sku").Pluck("sku", &sku)

	c.JSON(200, gin.H{"success": true, "sku": sku})
}
