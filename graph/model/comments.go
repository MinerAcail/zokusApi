package model
import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID              uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID          string  ` json:"userId"`
	ProductID       string  ` json:"productId"`
	Content         string      `json:"content"`
	ParentCommentID *uuid.UUID  `gorm:"type:uuid" json:"parentCommentId,omitempty"`  // Use UUID for parent comment
	Replies         []*Comment  `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`  // Foreign key points to ParentCommentID
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}
