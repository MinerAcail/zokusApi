package model

import (
	"time"

	"github.com/google/uuid"
)

// Cart represents the shopping cart for a user
// DeliveryStatus defines the status of the cart item
type DeliveryStatu string

const (
	//Ordered > Shipped  >  Arrived / Missing
	//update it

	StatusOrder   DeliveryStatu = "Ordered" // Item not yet delivered
	StatusShipped DeliveryStatu = "Shipped" // Item is Shipped so with the driver
	StatusArrived DeliveryStatu = "Arrived" // Item has arrived
	StatusMissing DeliveryStatu = "Missing" // Item is reported as missing
)

type PurchaseMetaDataStatus string

const (
	StatusPurchased PurchaseMetaDataStatus = "Purchased"
	StatusInCart    PurchaseMetaDataStatus = "inCart"
	StatusCancelled PurchaseMetaDataStatus = "Cancelled"
)

type ShippingMethod string

const (
	ShippingMethodPickup   ShippingMethod = "Pickup"
	ShippingMethodDelivery ShippingMethod = "Delivery"
)

type PaymentMethod string

const (
	PaymentMethodCash PaymentMethod = "Cash"
	PaymentMethodMomo PaymentMethod = "Momo"
)

type Cart struct {
	ID               uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserType         UserType            `gorm:"type:varchar(20)" json:"userType"` // Holds "User", "Vendor", or "SideSwapper"
	UserID           uuid.UUID           `gorm:"type:uuid" json:"userId"`
	UserName         string              `gorm:"-" json:"userName"`
	UserNumber       string              `gorm:"-" json:"userNumber"`
	Items            []*CartItem         `gorm:"foreignKey:CartID" json:"items,omitempty"`
	Location         []*DeliveryMetaData `gorm:"foreignKey:CartID" json:"locations,omitempty"`
	PurchaseMetaData []*PurchaseMetaData `gorm:"foreignKey:CartID" json:"purchaseMetaData,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

// CartItem represents an individual item in the cart
type CartItem struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Product         *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductID       uuid.UUID  `gorm:"type:uuid" json:"productId"`
	CartID          uuid.UUID  `gorm:"type:uuid" json:"cartId"`
	Quantity        int        `json:"quantity"`
	OrderID         *string    `json:"orderId,omitempty"`                          //this is the order id from the generatedid by GenerateOrderCartID
	OrderCartItemId *uuid.UUID `gorm:"type:uuid" json:"orderCartItemId,omitempty"` // Use this as the foreign key for the order
	// PurchaseMetaData *PurchaseMetaData `gorm:"foreignKey:CartItemID;constraint:OnDelete:CASCADE;" json:"purchaseMetaData,omitempty"`
	// DeliveryStatus   DeliveryStatus    `gorm:"foreignKey:CartItemID;constraint:OnDelete:CASCADE;" json:"deliveryStatus,omitempty"`
	// DeliveryMetaData DeliveryMetaData  `gorm:"foreignKey:CartItemID;constraint:OnDelete:CASCADE;" json:"deliveryMetaData,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Order struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CustomerID            uuid.UUID `gorm:"type:uuid" json:"customerId"`
	CustomerName          string    `json:"customerName"`
	Date                  string    `json:"date"`
	DeliveryAddress       string    `json:"deliveryAddress"`
	Status                string    `json:"status"`
	OrderItems            []*Group  `gorm:"foreignKey:OrderItemsID" json:"orderItems"`
	LocationName          *string   `json:"locationName,omitempty"`
	MomoNumber            *string   `json:"momoNumber,omitempty"`
	Latitude              *string   `json:"latitude,omitempty"`
	Longitude             *string   `json:"longitude,omitempty"`
	PaymentStatus         string    `json:"paymentStatus"`
	PaymentMethod         string    `json:"paymentMethod"`
	ShippingMethod        string    `json:"shippingMethod"`
	TrackingNumber        *string   `json:"trackingNumber,omitempty"`
	EstimatedDeliveryDate string    `json:"estimatedDeliveryDate"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// Group struct represents an item within a group, belonging to a CartItemGroup and a CartItem
type Group struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CartItem        *CartItem      `json:"cartItem" gorm:"foreignKey:CartItemID;constraint:OnDelete:CASCADE;"`
	CartItemID      uuid.UUID      `gorm:"type:uuid" json:"cartItemId"`
	CartItemGroup   *CartItemGroup `json:"cartItemGroup" gorm:"foreignKey:CartItemGroupID;constraint:OnDelete:CASCADE;"` // Add reference to CartItemGroup
	CartItemGroupID uuid.UUID      `gorm:"type:uuid" json:"cartItemGroupId"`                                             // Foreign key for CartItemGroup
	OrderItemsID    *uuid.UUID     `gorm:"type:uuid" json:"orderItemsId"`
	VendorID        *string        `json:"vendorId,omitempty"`
	QuantityGroup   int            `json:"quantityGroup"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

// CartItemGroup represents a collection of Group items that belong to the same CartItemGroup
type CartItemGroup struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GroupLabel string    `json:"groupLabel"`
	Items      []*Group  `json:"items" gorm:"foreignKey:CartItemGroupID;constraint:OnDelete:CASCADE;"` // Reference to related Groups
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type PurchaseMetaData struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Status     string    `json:"status"`
	CartItemID uuid.UUID `gorm:"type:uuid" json:"cartItemId"`
	CartID     uuid.UUID `gorm:"type:uuid" json:"cartId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type DeliveryMetaData struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	LocationName string    `json:"locationName"`
	MomoNumber   string    `json:"momoNumber"`
	CartItemID   uuid.UUID `gorm:"type:uuid" json:"cartItemId"`
	OrderID      string    `json:"orderId"`
	Latitude     string    `json:"latitude"`
	Longitude    string    `json:"longitude"`
	CartID       uuid.UUID `gorm:"type:uuid" json:"cartId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type DeliveryStatus struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Status      string    `json:"status"`
	CartItemID  uuid.UUID `gorm:"type:uuid" json:"cartItemId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PurchaseHistory struct {
	ID             string                 `json:"id"`
	Date           string                 `json:"date"`
	Items          []*PurchaseHistoryItem `json:"items"`
	Status         string                 `json:"status"`
	TrackingNumber *string                `json:"trackingNumber,omitempty"`
}

type PurchaseHistoryItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	ImageURL string `json:"imageUrl"`
}

type OrderItem struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Price    float64   `json:"price"`
}
