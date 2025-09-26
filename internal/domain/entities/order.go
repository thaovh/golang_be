package entities

import "github.com/google/uuid"

// Order represents an order entity in the domain
// Maps to BMSF_ORDER table in Oracle database
type Order struct {
	BaseEntity
	OrderNumber string      `json:"order_number" gorm:"size:50;uniqueIndex;not null" db:"ORDER_NUMBER"` // Maps to BMSF_ORDER.ORDER_NUMBER
	UserID      uuid.UUID   `json:"user_id" gorm:"type:varchar(36);not null;index" db:"USER_ID"`        // Maps to BMSF_ORDER.USER_ID
	TotalAmount float64     `json:"total_amount" gorm:"type:number(10,2);default:0" db:"TOTAL_AMOUNT"`  // Maps to BMSF_ORDER.TOTAL_AMOUNT
	Status      OrderStatus `json:"status" gorm:"size:20;default:'PENDING';not null" db:"STATUS"`       // Maps to BMSF_ORDER.STATUS
	Notes       string      `json:"notes" gorm:"size:1000" db:"NOTES"`                                  // Maps to BMSF_ORDER.NOTES
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// IsValid checks if the order status is valid
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// NewOrder creates a new order entity
func NewOrder(orderNumber string, userID uuid.UUID, totalAmount float64) *Order {
	order := &Order{
		BaseEntity:  NewBaseEntity(),
		OrderNumber: orderNumber,
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      OrderStatusPending,
	}
	return order
}

// Confirm confirms the order
func (o *Order) Confirm(updatedBy *uuid.UUID) {
	o.Status = OrderStatusConfirmed
	o.UpdateVersion(updatedBy)
}

// Ship ships the order
func (o *Order) Ship(updatedBy *uuid.UUID) {
	o.Status = OrderStatusShipped
	o.UpdateVersion(updatedBy)
}

// Deliver delivers the order
func (o *Order) Deliver(updatedBy *uuid.UUID) {
	o.Status = OrderStatusDelivered
	o.UpdateVersion(updatedBy)
}

// Cancel cancels the order
func (o *Order) Cancel(updatedBy *uuid.UUID) {
	o.Status = OrderStatusCancelled
	o.UpdateVersion(updatedBy)
}
