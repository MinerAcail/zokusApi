package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Product struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name              *string        `json:"name,omitempty"`
	Price             *string        `json:"price,omitempty"`
	State             *string        `json:"state,omitempty"`
	Title             *string        `json:"title,omitempty"`
	Brand             *string        `json:"brand,omitempty"`
	ImageUrls         pq.StringArray `gorm:"type:text[]" json:"imageUrls,omitempty"`
	VendorsBucketName *string        `json:"vendorsBucketName,omitempty"`
	Vendor            *Vendor        `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	VendorID          uuid.UUID      `gorm:"type:uuid" json:"vendorId,omitempty"`
	CreatedAt         *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time     `json:"updatedAt,omitempty"`
}

type ProductRequest struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key;" json:"id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserID      uuid.UUID      `gorm:"type:uuid" json:"userId"`
	Vendor      *Vendor        `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	VendorID    uuid.UUID      `gorm:"type:uuid" json:"vendorId"`
	Description *string        `json:"description,omitempty"`
	ImageUrls   pq.StringArray `gorm:"type:text[]" json:"imageUrls,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type FavoriteProduct struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserType  UserType  `gorm:"type:varchar(20)" json:"userType"` // Holds "User", "Vendor", or "SideSwapper"
	UserID    uuid.UUID `gorm:"type:uuid" json:"userId"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"productId"`
	CreatedAt time.Time `json:"createdAt"`
}
