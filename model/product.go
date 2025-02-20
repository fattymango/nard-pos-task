package model

type Product struct {
	ID int `gorm:"primaryKey,autoIncrement" json:"id"`
	// ForeignKey
	TenantID int `gorm:"column:tenant_id" json:"tenant_id"`
	// Info
	Name  string  `gorm:"column:name" json:"name"`
	Price float64 `gorm:"column:price" json:"price"`

	CommonModel
}

type ProductSales struct {
	ProductID  int32   `json:"product_id"`
	TotalSales float64 `json:"total_sales"`
}
