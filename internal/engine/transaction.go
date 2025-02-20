package engine

import (
	"fmt"
	"multitenant/model"
	"time"
)

const (
	maxRetries     = 3
	initialBackoff = 100 * time.Millisecond
	backoffFactor  = 2
)

// CreateTransaction queues a transaction for processing
func (e *Engine) CreateTransaction(tx *model.Transaction) error {

	ok, err := e.TanentExists(tx.TenantID)
	if err != nil {
		return fmt.Errorf("failed to check tenant existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("tenant %d does not exist", tx.TenantID)
	}

	ok, err = e.ProductExists(tx.ProductID)
	if err != nil {
		return fmt.Errorf("failed to check product existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("product %d does not exist", tx.ProductID)
	}

	ok, err = e.BranchExists(tx.BranchID)
	if err != nil {
		return fmt.Errorf("failed to check branch existence: %w", err)
	}

	if !ok {
		return fmt.Errorf("branch %d does not exist", tx.BranchID)
	}

	// e.logger.Debugf("Adding transaction to queue: %v", tx)
	e.txQueue <- tx
	return nil
}
