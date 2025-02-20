package app

import (
	"fmt"
	"multitenant/model"
	"multitenant/pkg/config"
	"multitenant/pkg/db"
)

func Migrate(cfg *config.Config, db *db.DB) error {
	err := db.AutoMigrate(&model.Tenant{}, &model.Branch{}, &model.Product{}, &model.Transaction{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %s", err)
	}

	return nil
}
