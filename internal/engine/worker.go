package engine

import (
	"multitenant/model"
	"multitenant/pkg/metrics"
	"time"
)

func (e *Engine) startWorkerPool(n int) {
	for i := 0; i < n; i++ {
		e.wg.Add(1)
		go e.worker(i)
	}
}

func (e *Engine) worker(id int) {
	defer e.wg.Done()
	for {
		select {
		case batch := <-e.workerQueue:
			e.logger.Debugf("Worker %d processing batch of %d transactions", id, len(batch))
			e.flushBatch(batch)
		case <-e.ctx.Done():
			e.logger.Debugf("Worker %d stopping", id)
			return
		}
	}
}

func (e *Engine) processBatches() {
	e.logger.Debugf("Starting transaction processor")
	defer e.wg.Done()

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			e.logger.Debugf("Stopping transaction processor")
			if len(e.batchBuffer) > 0 {
				batch := e.drainBatch()
				e.workerQueue <- batch // Add remaining batch to the queue
			}
			return

		case tx := <-e.txQueue:
			e.mu.Lock()
			e.refreshSellingProducts(tx.TenantID)
			e.batchBuffer = append(e.batchBuffer, tx)
			if len(e.batchBuffer) >= batchSize {
				e.logger.Debugf("Flushing batch due to size")
				batch := e.drainBatch()
				e.workerQueue <- batch // Add batch to the queue for workers
			}
			e.mu.Unlock()

		case <-ticker.C:
			// e.logger.Debugf("Flushing batch due to timeout")
			e.mu.Lock()
			if len(e.batchBuffer) > 0 {
				batch := e.drainBatch()
				e.workerQueue <- batch // Add batch to the queue for workers
			}
			e.mu.Unlock()
		}
	}
}

func (e *Engine) drainBatch() []*model.Transaction {
	batch := e.batchBuffer
	e.batchBuffer = nil // Reset buffer
	return batch
}

// flushBatch flushes the transaction batch to the database and updates the cache
func (e *Engine) flushBatch(batch []*model.Transaction) {
	if len(batch) == 0 {
		return
	}

	e.logger.Debugf("Flushing batch of %d transactions", len(batch))

	query := "INSERT INTO transaction (tenant_id, branch_id, product_id, quantity_sold, price_per_unit) VALUES "
	values := []interface{}{}

	salesCount := make(map[int32]int32)
	for _, tx := range batch {
		query += "(?, ?, ?, ?, ?),"
		values = append(values, tx.TenantID, tx.BranchID, tx.ProductID, tx.QuantitySold, tx.PricePerUnit)
		salesCount[tx.ProductID] += tx.QuantitySold
		e.metrics.IncrementTransactionsProcessed(metrics.TX_PROCCESSED_STATUS_SUCCESS)
	}
	query = query[:len(query)-1] // Remove trailing comma

	var err error
	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = e.db.Exec(query, values...).Error
		if err == nil {
			e.logger.Debugf("Successfully inserted %d transactions", len(batch))
			break
		}
		e.logger.Errorf("Batch insert failed (attempt %d/%d): %v", attempt+1, maxRetries, err)
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= backoffFactor
		}
	}

	if err != nil {
		e.logger.Errorf("Batch insert ultimately failed after %d attempts: %v", maxRetries, err)
		return
	}

	// Update cache for total sales
	for productID, totalSaleAmount := range salesCount {
		e.CacheTotalSalesPerProduct(batch[0].TenantID, productID, float64(totalSaleAmount))
	}
}

func (e *Engine) refreshSellingProducts(tenantID int32) error {
	var t int
	n, ok := e.mp.Load(topSellingProductsKey(tenantID))
	if !ok {
		e.mp.Store(topSellingProductsKey(tenantID), 0)
		n = 0
	}
	t = n.(int) + 1
	e.mp.Store(topSellingProductsKey(tenantID), n.(int)+1)
	if n.(int)+1 >= autoRefreshTxThreshold {
		e.logger.Debugf("Refreshing cache after %d transactions", n.(int)+1)
	}
	if t >= autoRefreshTxThreshold {
		e.logger.Debugf("Refreshing cache after %d transactions", t)
		e.mp.Store(topSellingProductsKey(tenantID), 0)
		e.logger.Debugf("Refreshing cache after %d transactions", t)
		_, err := e.GetTopSellingProducts(tenantID)
		if err != nil {
			e.logger.Errorf("Failed to refresh cache: %v", err)
		}
	}

	return nil
}
