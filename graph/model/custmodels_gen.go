
package model

import (
	"time"
	"github.com/lib/pq"
	"github.com/google/uuid"

)



type User struct {
	ID   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string     `json:"name"`
	Number    string     `json:"number"`
	Pin       *string    `json:"pin,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
type Vender struct {
	ID   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string        `json:"name"`
	Number      string        `json:"number"`
	Pin         *string       `json:"pin,omitempty"`
	Product     []*Product  `gorm:"foreignKey:VenderID" json:"product,omitempty"`
	IsCompleted *bool         `json:"isCompleted,omitempty"`
	CreatedAt   *time.Time    `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time    `json:"updatedAt,omitempty"`
}

type Product struct {
    ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    Name      string          `json:"name"`
    Price     string          `json:"price"`
    State     string          `json:"state"`
    Title     string          `json:"title"`
    Brand     string          `json:"brand"`
    ImageUrls pq.StringArray  `gorm:"type:text[]" json:"imageUrls,omitempty"`
    VenderID  uuid.UUID       `gorm:"type:uuid" json:"venderId"`
}



type UserOrder struct {
    ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
    User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
    UserID      uuid.UUID      `gorm:"type:uuid" json:"userId"`
    // Products    []*Product  `gorm:"many2many:user_order_products;" json:"products,omitempty"`
    Products    []*Product  `gorm:"many2many:user_order_products;" json:"products,omitempty"`

    Quantity    string         `json:"quantity"`
    IsCompleted *bool          `json:"isCompleted,omitempty"`
    CreatedAt   *time.Time     `json:"createdAt,omitempty"`
    UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
}




type VenderWithItProduct struct {
	ID   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	User      *Vender       `json:"user,omitempty"`
	Product   []*Product `json:"product,omitempty"`
	CreatedAt *time.Time    `json:"createdAt,omitempty"`
	UpdatedAt *time.Time    `json:"updatedAt,omitempty"`
}



type ProductRequest struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
	User        *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserID      uuid.UUID     `gorm:"type:uuid" json:"userId"`
	Vender      *Vender       `gorm:"foreignKey:VenderID" json:"vender,omitempty"`
	VenderID    uuid.UUID     `gorm:"type:uuid" json:"venderId"`
	Description *string       `json:"description,omitempty"`
    ImageUrls pq.StringArray  `gorm:"type:text[]" json:"imageUrls,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}
