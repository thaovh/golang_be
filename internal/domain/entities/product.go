package entities

import "github.com/google/uuid"

// Product represents a product entity in the domain
// Maps to BMSF_PRODUCT table in Oracle database
type Product struct {
	BaseEntity
	Code        string        `json:"code" gorm:"size:50;uniqueIndex;not null" db:"CODE"`                 // Maps to BMSF_PRODUCT.CODE
	Name        string        `json:"name" gorm:"size:255;not null" db:"NAME"`                 // Maps to BMSF_PRODUCT.NAME
	Description string        `json:"description" gorm:"size:1000" db:"DESCRIPTION"`   // Maps to BMSF_PRODUCT.DESCRIPTION
	Price       float64       `json:"price" gorm:"type:number(10,2);default:0" db:"PRICE"`               // Maps to BMSF_PRODUCT.PRICE
	Category    string        `json:"category" gorm:"size:100;index" db:"CATEGORY"`         // Maps to BMSF_PRODUCT.CATEGORY
	Status      ProductStatus `json:"status" gorm:"size:20;default:'DRAFT';not null" db:"STATUS"`           // Maps to BMSF_PRODUCT.STATUS
}

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
	ProductStatusDraft    ProductStatus = "DRAFT"
	ProductStatusArchived ProductStatus = "ARCHIVED"
)

// IsValid checks if the product status is valid
func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusActive, ProductStatusInactive, ProductStatusDraft, ProductStatusArchived:
		return true
	default:
		return false
	}
}

// NewProduct creates a new product entity
func NewProduct(code, name, description, category string, price float64) *Product {
	product := &Product{
		BaseEntity:  NewBaseEntity(),
		Code:        code,
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		Status:      ProductStatusDraft,
	}
	return product
}

// Activate activates the product
func (p *Product) Activate(updatedBy *uuid.UUID) {
	p.Status = ProductStatusActive
	p.UpdateVersion(updatedBy)
}

// Deactivate deactivates the product
func (p *Product) Deactivate(updatedBy *uuid.UUID) {
	p.Status = ProductStatusInactive
	p.UpdateVersion(updatedBy)
}

// Archive archives the product
func (p *Product) Archive(updatedBy *uuid.UUID) {
	p.Status = ProductStatusArchived
	p.UpdateVersion(updatedBy)
}

// UpdatePrice updates the product price
func (p *Product) UpdatePrice(price float64, updatedBy *uuid.UUID) {
	p.Price = price
	p.UpdateVersion(updatedBy)
}
