package types

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement:false;`
	ProductID uint   `json:"id" gorm:"primaryKey;column:product_id;type:bigint;autoIncrement:true;not null"`
	SKU       string `json:"sku" gorm:"type:varchar(100)"`
	// Virtual       bool    `json:"virtual"`
	// Downloadable  bool    `json:"downloadable"`
	Price         float32 `json:"price"`
	OnSale        bool    `json:"onSale" gorm:"column:onsale"`
	StockQuantity float32 `json:"stockQuantity"`
	StockStatus   string  `json:"stockStatus" gorm:"type:varchar(100);default:'instock'"`
	// RatingCount   int     `json:"ratingCount" gorm:"default:0"`
	// AverageRating float32 `json:"averageRating" gorm:"default:0.00"`
	// TotalSales    int     `json:"totalSales" gorm:"default:0"`
	// TaxStatus     string  `json:"taxStatus" gorm:"default:taxable"`
	// TaxClass      string  `json:"taxClass"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Gender    string `json:"gender"`
	Season    string `json:"season"`
	Location  string `json:"location"`
	Image     []byte `json:"image"`
	Timestamp int64  `json:"timestamp"`
}

type Item2 struct {
	gorm.Model
	Name     string  `json:"name"`
	SKU      string  `json:"name"`
	Category string  `json:"category"` //Gender and such
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Size     string  `json:"size"`
}

func (Item) TableName() string {
	return "AP6_wc_product_meta_lookup"
}

func (Item2) TableName() string {
	return "Item"
}
