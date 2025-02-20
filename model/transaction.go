package model

type TransactionStatus uint8

const (
	TransactionStatusPending TransactionStatus = iota + 1
	TransactionStatusCompleted
	TransactionStatusFailed
)

type Transaction struct {
	ID int32 `gorm:"primaryKey,autoIncrement" json:"id"`
	// ForeignKey
	TenantID  int32 `gorm:"column:tenant_id" json:"tenant_id"`
	BranchID  int32 `gorm:"column:branch_id" json:"branch_id"`
	ProductID int32 `gorm:"column:product_id" json:"product_id"`
	// Info
	QuantitySold int32   `gorm:"column:quantity_sold" json:"quantity_sold"`
	PricePerUnit float64 `gorm:"column:price_per_unit" json:"price_per_unit"`

	// Control
	Status TransactionStatus `gorm:"column:status" json:"status"`
	CommonModel
}

type CrtTransaction struct {
	TenantID     int32   `json:"tenant_id" validate:"required"`
	BranchID     int32   `json:"branch_id" validate:"required"`
	ProductID    int32   `json:"product_id" validate:"required"`
	QuantitySold int32   `json:"quantity_sold" validate:"required"`
	PricePerUnit float64 `json:"price_per_unit" validate:"required"`
}

func (crt *CrtTransaction) ToTransaction() *Transaction {
	return &Transaction{
		TenantID:     crt.TenantID,
		BranchID:     crt.BranchID,
		ProductID:    crt.ProductID,
		QuantitySold: crt.QuantitySold,
		PricePerUnit: crt.PricePerUnit,
		Status:       TransactionStatusPending,
	}
}
