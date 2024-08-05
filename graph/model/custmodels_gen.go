package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string     `json:"name"`
	Type      *string    `json:"type,omitempty"`
	Number    string     `json:"number"`
	Pin       *string    `json:"pin,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
type Vendor struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string     `json:"name"`
	Type        *string    `json:"type,omitempty"`
	Number      string     `json:"number"`
	Pin         *string    `json:"pin,omitempty"`
	Product     []*Product `gorm:"foreignKey:VendorID" json:"product,omitempty"`
	FileUrls    []*FileURL `json:"fileUrls,omitempty"`
	IsCompleted *bool      `json:"isCompleted,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type FileURL struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	VendorID  uuid.UUID  `gorm:"type:uuid" json:"vendorId"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type Product struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
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
type UserOrder struct {
	ID       uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	User     *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserID   uuid.UUID  `gorm:"type:uuid" json:"userId"`
	Products []*Product `gorm:"many2many:user_order_products;" json:"products,omitempty"`

	Quantity    string     `json:"quantity"`
	IsCompleted *bool      `json:"isCompleted,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type VendorWithItProduct struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	User      *Vendor    `json:"user,omitempty"`
	Product   []*Product `json:"product,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type ProductRequest struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserID      uuid.UUID      `gorm:"type:uuid" json:"userId"`
	Vendor      *Vendor        `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	VendorID    uuid.UUID      `gorm:"type:uuid" json:"vendorId"`
	Description *string        `json:"description,omitempty"`
	ImageUrls   pq.StringArray `gorm:"type:text[]" json:"imageUrls,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

func (User) IsReturns()   {}
func (Vendor) IsReturns() {}
