package app

import (
	"fmt"
	"multitenant/model"
	"multitenant/pkg/db"
)

func Seed(db *db.DB) error {
	// Seed data
	tenants := []model.Tenant{
		{
			ID:   1,
			Name: "Tenant 1",
		},
		{
			ID:   2,
			Name: "Tenant 2",
		},
		{
			ID:   3,
			Name: "Tenant 3",
		},
	}

	for i, tenant := range tenants {
		err := db.Create(&tenant).Error
		if err != nil {
			return fmt.Errorf("failed to seed tenant: %w", err)
		}

		branches := []model.Branch{
			{
				ID:       (i * 5) + 1,
				TenantID: tenant.ID,
				Name:     "Branch 1",
			},
			{
				ID:       (i * 5) + 2,
				TenantID: tenant.ID,
				Name:     "Branch 2",
			},
			{
				ID:       (i * 5) + 3,
				TenantID: tenant.ID,
				Name:     "Branch 3",
			},
			{
				ID:       (i * 5) + 4,
				TenantID: tenant.ID,
				Name:     "Branch 4",
			},
			{
				ID:       (i * 5) + 5,
				TenantID: tenant.ID,
				Name:     "Branch 5",
			},
		}

		err = db.Create(&branches).Error
		if err != nil {
			return fmt.Errorf("failed to seed branch: %w", err)
		}

		products := []model.Product{
			{
				ID:       (i * 3) + 1,
				TenantID: tenant.ID,
				Name:     fmt.Sprintf("%s - Product 1", tenant.Name),
				Price:    1000,
			},
			{
				ID:       (i * 3) + 2,
				TenantID: tenant.ID,
				Name:     fmt.Sprintf("%s - Product 2", tenant.Name),
				Price:    2000,
			},
			{
				ID:       (i * 3) + 3,
				TenantID: tenant.ID,
				Name:     fmt.Sprintf("%s - Product 3", tenant.Name),
				Price:    3000,
			},
		}

		err = db.Create(&products).Error
		if err != nil {
			return fmt.Errorf("failed to seed product: %w", err)
		}

	}

	return nil

}
