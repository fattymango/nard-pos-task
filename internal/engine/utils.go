package engine

import (
	"fmt"
	"multitenant/model"
)

func (e *Engine) TanentExists(tenantID int32) (bool, error) {
	// check map first
	_, ok := e.mp.Load(tenantKey(tenantID))
	if ok {
		return true, nil
	}
	var count int64
	err := e.db.Model(&model.Tenant{}).Where("id = ?", tenantID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check tenant existence: %w", err)
	}

	if count > 0 {
		e.mp.Store(tenantKey(tenantID), struct{}{})
	}

	return count > 0, nil
}

func (e *Engine) ProductExists(productID int32) (bool, error) {
	// check map first
	_, ok := e.mp.Load(productKey(productID))
	if ok {
		return true, nil
	}

	var count int64
	err := e.db.Model(&model.Product{}).Where("id = ?", productID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}

	if count > 0 {
		e.mp.Store(productKey(productID), struct{}{})
	}

	return count > 0, nil
}

func (e *Engine) BranchExists(branchID int32) (bool, error) {
	// check map first
	_, ok := e.mp.Load(branchKey(branchID))
	if ok {
		return true, nil
	}

	var count int64
	err := e.db.Model(&model.Branch{}).Where("id = ?", branchID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check branch existence: %w", err)
	}

	if count > 0 {
		e.mp.Store(branchKey(branchID), struct{}{})
	}

	return count > 0, nil
}
