package model

import (
	"time"

	"github.com/google/uuid"
)

type FileURL struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	VendorID  uuid.UUID  `gorm:"type:uuid" json:"vendorId"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (User) IsReturns()   {}
func (Vendor) IsReturns() {}
