package model

import (
	"time"

	"github.com/google/uuid"
)

type SideSwapper struct {
	ID                   uuid.UUID                   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name                 string                      `json:"name"`
	Type                 *string                     `json:"type,omitempty"`
	Number               string                      `json:"number"`
	Pin                  *string                     `json:"pin,omitempty"`
	ProductsIds          []*uuid.UUID                `gorm:"-" json:"productsIds"`
	Products             []*SideSwapperProduct       `gorm:"foreignKey:SideSwapperID" json:"products"`
	InventoryMetaDataIds []*uuid.UUID                `gorm:"-" json:"inventoryMetaDataIds"`
	InventoryMetaData    []*ProductInventoryMetaData `gorm:"foreignKey:SideSwapperID" json:"inventoryMetaData"`
	PricingID            *uuid.UUID                  `gorm:"type:uuid" json:"pricingId"` // Reference to Pricing
	Pricing              *Pricing                    `gorm:"foreignKey:PricingID" json:"pricing"`
	CreatedAt            *time.Time                  `json:"createdAt,omitempty"`
	UpdatedAt            *time.Time                  `json:"updatedAt,omitempty"`
}
type Pricing struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`

	BasePrice        float64 `json:"basePrice"`
	MarkupPercentage float64 `json:"markupPercentage"`
}

type SideSwapperProduct struct {
	ID            uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	SideSwapper   *SideSwapper `json:"sideSwapper,omitempty"`
	SideSwapperID uuid.UUID    `gorm:"type:uuid" json:"sideSwapperID"`
	Vendor        *Vendor      `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	VendorID      uuid.UUID    `gorm:"type:uuid" json:"vendorId"`
	Product       *Product     `json:"product,omitempty"`
	ProductID     uuid.UUID    `gorm:"type:uuid" json:"productId"`
	PricingID     *uuid.UUID   `gorm:"type:uuid" json:"pricingId"` // Reference to Pricing
	Pricing       *Pricing     `gorm:"foreignKey:PricingID" json:"pricing"`
}


type ProductInventoryMetaData struct {
	ID              uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	SideSwapper     *SideSwapper `gorm:"foreignKey:SideSwapperID" json:"sideSwapper"`
	SideSwapperID   uuid.UUID    `gorm:"type:uuid" json:"sideSwapperID"`
	StockLevel      int          `json:"stockLevel"`
	PercentageLevel float64      `json:"PercentageLevel"`
	UpdatedAt       *time.Time   `json:"updatedAt,omitempty"`
}
