package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name      string     `json:"name"`
	Type      *string    `json:"type,omitempty"`
	Number    string     `json:"number"`
	Pin       *string    `json:"pin,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
type UserOrder struct {
	ID       uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key;primaryKey" json:"id"`
	User     *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserID   uuid.UUID  `gorm:"type:uuid" json:"userId"`
	Products []*Product `gorm:"many2many:user_order_products;" json:"products,omitempty"`

	Quantity    string     `json:"quantity"`
	IsCompleted *bool      `json:"isCompleted,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}
