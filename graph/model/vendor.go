package model

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

type Vendor struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name          string     `json:"name"`
	Type          *string    `json:"type,omitempty"`
	Number        string     `json:"number"`
	Pin           *string    `json:"pin,omitempty"`
	Product       []*Product `gorm:"foreignKey:VendorID" json:"product,omitempty"`
	FileUrls      []*FileURL `json:"fileUrls,omitempty"`
	ProfileURL    *string    `json:"profileUrl,omitempty"`
	BackgroundURL *string    `json:"backgroundUrl,omitempty"`
	IsCompleted   *bool      `json:"isCompleted,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

type VendorWithItProduct struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	User      *Vendor    `json:"user,omitempty"`
	Product   []*Product `json:"product,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type UploadVendorMateData struct {
	Profile    *graphql.Upload `json:"profile,omitempty"`
	Background *graphql.Upload `json:"background,omitempty"`
}
type NewVendor struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	Pin    string `json:"pin"`
}


type FavoriteVendor struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
    UserID    uuid.UUID `gorm:"type:uuid" json:"userId"`
    Vendor    *Vendor   `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
    VendorID  uuid.UUID `gorm:"type:uuid" json:"vendorId"`
    CreatedAt time.Time `json:"createdAt"`
}

type FavoriteVendorMateData struct {
	FavoriteVendors     []*FavoriteVendor `json:"favoriteVendors"`
	Count               int               `json:"count"`
	MostFavoritedVendor *Vendor           `json:"mostFavoritedVendor,omitempty"`
}

