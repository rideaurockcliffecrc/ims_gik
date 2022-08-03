package items_temp

import (
	"GIK_Web/database"
	"GIK_Web/utils"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type post struct {
	ID       uint   `json:"id"`
	PostName string `json:"name"`
}

type foundItem struct {
	Name      string `json:"name"`
	ProductID uint   `json:"productId"`
}

func GetIdsInBulk(c *gin.Context) {
	// get plaintext body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	// for every line, in the body, append to slice
	var productNames []string
	productNames = append(productNames, strings.Split(string(body), "\n")...)

	failedProductNames := []string{}

	// get ids from product names
	var productIds []foundItem
	for _, productName := range productNames {

		if len(productName) == 0 {
			continue
		}

		originalProductName := productName

		productName = utils.StandardizeSpaces(productName)

		combinations := []string{}

		giveUp := false

		combinations = append(combinations, productName)

		combinations = append(combinations, strings.ReplaceAll(productName, "'", ""))

		splitString := strings.Split(strings.ToLower(productName), " ")

		// if !strings.Contains(splitString[0], "'") {
		// 	splitString2 := splitString
		// 	splitString2[0] += "'s"
		// 	combinations = append(combinations, strings.Join(splitString2, " "))
		// }

		if !strings.Contains(strings.ToLower(productName), "all season") {
			// split
			splitString = append(splitString, "")
			copy(splitString[2:], splitString[1:])
			splitString[1] = "all season"
			combinations = append(combinations, strings.Join(splitString, " "))
			fmt.Println(splitString)
		}

		// make first word not plural
		if productName[len(productName)-1:] == "s" {
			productName = strings.ReplaceAll(productName, "s", "")
			combinations = append(combinations, productName)
		}

		post := post{}
		for _, combination := range combinations {
			if giveUp {
				break
			}

			database.Database.Raw("SELECT * FROM `AP6_posts` WHERE `post_title` like ?", fmt.Sprintf("%%%s%%", combination)).Scan(&post)
			if post.PostName != "" {
				productIds = append(productIds, foundItem{
					ProductID: post.ID,
					Name:      post.PostName,
				})
				giveUp = true
			}
		}

		if post.ID == 0 {
			failedProductNames = append(failedProductNames, originalProductName)
		}

	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"successful": productIds,
			"failed":     failedProductNames,
		},
	})
}

func LookupItem(c *gin.Context) {
	// product id
	productId := c.Query("product_id")

	if productId == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Product ID not provided",
		})
		return
	}

	// convert to int
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	var postData post
	database.Database.Raw("SELECT * FROM `AP6_posts` WHERE `ID` = ?", productIdInt).Scan(&postData)

	c.JSON(200, gin.H{
		"success": true,
		"data":    postData,
	})

}
