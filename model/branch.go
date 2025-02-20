package model

type Branch struct {
	ID int `gorm:"primaryKey,autoIncrement" json:"id"`
	// ForeignKey
	TenantID int `gorm:"column:tenant_id" json:"tenant_id"`
	// Info
	Name string `gorm:"column:name" json:"name"`

	CommonModel
}
