package model

type Tenant struct {
	ID int `gorm:"primaryKey,autoIncrement" json:"id"`
	// Info
	Name string `gorm:"column:name" json:"name"`

	CommonModel
}
