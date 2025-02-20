package engine

import (
	"fmt"
	"multitenant/model"
)

func (e *Engine) GetTotalSalesPerProduct(tenantID, productID int32) (float64, error) {
	var totalSales float64
	totalSales, err := e.GetCachedTotalSalesPerProduct(tenantID, productID)
	if err != nil && !isCacheMiss(err) {
		return 0, fmt.Errorf("failed to get total sales for tenant %d, product %d: %s", tenantID, productID, err)
	}

	if totalSales > 0 {
		return totalSales, nil
	}

	e.metrics.IncrementCacheMisses()
	// safe net, retrieve from db

	err = e.db.Raw(`
		SELECT COALESCE(SUM(quantity_sold * price_per_unit), 0) 
		FROM transaction
		WHERE tenant_id = ? AND product_id = ?;
	`, tenantID, productID).Scan(&totalSales).Error

	if err != nil {
		return 0, fmt.Errorf("failed to retrieve total sales from DB: %w", err)
	}

	err = e.CacheTotalSalesPerProduct(tenantID, productID, totalSales)
	if err != nil {
		e.logger.Errorf("Failed to cache total sales for tenant %d, product %d: %v", tenantID, productID, err)
	}
	return totalSales, nil
}

func (e *Engine) GetTopSellingProducts(tenantID int32) ([]model.ProductSales, error) {
	var topProducts []model.ProductSales

	topProducts, err := e.GetCachedTopSellingProducts(tenantID)
	if err != nil && !isCacheMiss(err) {
		return nil, fmt.Errorf("failed to get top selling products from cache: %s", err)
	}

	if len(topProducts) > 0 {
		return topProducts, nil
	}

	e.metrics.IncrementCacheMisses()

	// safe net, retrieve from db
	err = e.db.Raw(`
		SELECT p.id as product_id, COALESCE(SUM(t.quantity_sold * t.price_per_unit), 0) as total_sales
		FROM product p
		LEFT JOIN transaction t ON p.id = t.product_id AND t.tenant_id = ?
		WHERE p.tenant_id = ?
		GROUP BY p.id
		ORDER BY total_sales DESC
		LIMIT ?;
	`, tenantID, tenantID, 10).Scan(&topProducts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve top selling products from DB: %w", err)
	}

	err = e.CacheTopSellingProducts(tenantID, topProducts)
	if err != nil {
		e.logger.Errorf("Failed to cache top selling products: %v", err)
	}

	return topProducts, nil
}
