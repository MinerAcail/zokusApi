package model

import (
	"time"

	"github.com/google/uuid"
)

// Cart represents the shopping cart for a user
type Cart struct {
	ID        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserType  UserType    `gorm:"type:varchar(20)" json:"userType"` // Holds "User", "Vendor", or "SideSwapper"
	UserID    uuid.UUID   `gorm:"type:uuid" json:"userId"`
	Items     []*CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// CartItem represents an individual item in the cart
type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"productId"`
	CartID    uuid.UUID `gorm:"type:uuid" json:"cartId"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
