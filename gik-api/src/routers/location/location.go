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
	ItemID int    `json:"itemId"`
}

type lookupData struct {
	location
	Item types.Item `json:"product"`
}

type listData struct {
	location
	Item types.Item `json:"product"`
}

func ListLocation(c *gin.Context) {
	name := c.Query("name")
	letter := c.Query("letter")
	product := c.Query("product")

	locations := []location{}

	baseQuery := database.Database.Model(&types.Location{})

	if name != "" {
		baseQuery = baseQuery.Where("name = ?", name)
	}

	if letter != "" {
		baseQuery = baseQuery.Where("letter = ?", letter)
	}

	if product != "" {
		// check if it's an integer
		if _, err := strconv.Atoi(product); err == nil {
			// it's an integer
			baseQuery = baseQuery.Where("item_id = ?", product)
		} else {
			// it's a product name
			baseQuery = baseQuery.Where("item_id IN (SELECT product_id FROM AP6_wc_product_meta_lookup WHERE name LIKE ?)", fmt.Sprintf("%%%s%%", product))
		}
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
		err := database.Database.Model(&types.Item{}).Where("product_id = ?", location.ItemID).Scan(&item).Error

		if err != nil {
			continue
		}

		response = append(response, listData{
			location: location,
			Item:     item,
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

	database.Database.Model(&location{}).Where("name = ?", json.Name).Update("item_id", json.ItemID)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Updated location",
	})

	utils.CreateSimpleLog(c, "Updated location "+json.Name)

}

type addRequest struct {
	Name   string `json:"name" binding:"required"`
	Letter string `json:"letter" binding:"required"`
	Item   string `json:"item" binding:"required"`
}

func AddLocation(c *gin.Context) {
	json := addRequest{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid fields",
		})
		return
	}

	itemInt, err := strconv.Atoi(json.Item)

	// check if Item is an integer
	if err == nil {
		// it's an integer
		err := database.Database.Create(&types.Location{Name: json.Name, Letter: json.Letter, ItemID: itemInt}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Unable to create location",
			})
			return
		}
	} else {
		// it's a product name
		var item types.Item
		err := database.Database.Model(&types.Item{}).Where("name = ?", json.Item).Scan(&item).Error
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Unable to find by product name. Try the ID instead.",
			})
			return
		}
		err = database.Database.Create(&location{Name: json.Name, Letter: json.Letter, ItemID: int(item.ProductID)}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Unable to create location",
			})
			return
		}
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
		err := database.Database.Model(&types.Item{}).Where("product_id = ?", location.ItemID).Scan(&item).Error

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
	if err := database.Database.Model(&types.Item{}).Where("product_id = ?", location.ItemID).First(&item).Error; err != nil {
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
