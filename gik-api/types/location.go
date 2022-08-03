package types

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	Name   string `gorm:"name" json:"name"`
	Letter string `gorm:"type:char(1)" json:"letter"`
	ItemID int    `gorm:"index" json:"itemId"`
	Item   Item   `gorm:"foreignkey:ItemID;references:product_id" json:"product"`
}
