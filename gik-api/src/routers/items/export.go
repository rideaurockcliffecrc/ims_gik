package items

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
		ID       uint    `gorm:"primaryKey;autoIncrement:false" csv:"-"`
		Name     string  `json:"name"`
		SKU      string  `json:"name"`
		Category string  `json:"category"` //Gender and such
		Price    float32 `json:"price"`
		Quantity int     `json:"quantity"`
		Size     string  `json:"size"`
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
