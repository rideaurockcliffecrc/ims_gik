package items_temp

import (
	"GIK_Web/database"
	"GIK_Web/types"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"os"
)

func ExportItems(c *gin.Context) {
	type export struct {
		ID            uint    `gorm:"primaryKey;autoIncrement:false; csv:"-"`
		ProductID     uint    `json:"id" csv:"id" gorm:"primaryKey;column:product_id;type:bigint;autoIncrement:true;not null"`
		SKU           string  `json:"sku" csv:"sku" gorm:"type:varchar(100)"`
		Price         float32 `json:"price" csv:"price"`
		StockQuantity float32 `json:"stockQuantity" csv:"stockQuantity"`
		Name          string  `json:"name" csv:"name"`
		Category      string  `json:"category" csv:"category"`
		Gender        string  `json:"gender" csv:"gender"`
		Season        string  `json:"season" csv:"season"`
		Location      string  `json:"location" csv:"location"`
	}
	items := []export{}
	database.Database.Model(&types.Item{}).Find(&items)

	name := uuid.New().String() + ".csv"
	csvFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = gocsv.MarshalFile(&items, csvFile)
	c.File(name)

	os.Remove(name)
}
