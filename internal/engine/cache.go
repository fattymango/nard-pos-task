package engine

import (
	"encoding/json"
	"fmt"
	"multitenant/model"
	"time"
)

const (
	_Tenant_TTL               = int(10 * time.Minute)
	_PRODUCT_TTL              = int(10 * time.Minute)
	_BRANCH_TTL               = int(10 * time.Minute)
	_SALES_PER_PRODUCT_TTL    = int(10 * time.Minute)
	_TOP_SELLING_PRODUCTS_TTL = int(10 * time.Second)
)

var (
	tenantTotalSalesKey = func(tenantID int32, productID int32) string {
		return fmt.Sprintf("tenant:%d:product:%d:total_sales", tenantID, productID)
	}

	topSellingProductsKey = func(teantID int32) string {
		return fmt.Sprintf("tenant:%d:top_selling_products", teantID)
	}

	tenantKey = func(tenantID int32) string {
		return fmt.Sprintf("tenant:%d", tenantID)
	}

	productKey = func(productID int32) string {
		return fmt.Sprintf("product:%d", productID)
	}

	branchKey = func(branchID int32) string {
		return fmt.Sprintf("branch:%d", branchID)
	}

	isCacheMiss = func(err error) bool {
		return err.Error() == "redis: nil"
	}
)

func (e *Engine) CacheTenant(tenantID int32, tenant *model.Tenant) error {
	key := tenantKey(tenantID)
	data, err := json.Marshal(tenant)
	if err != nil {
		return fmt.Errorf("failed to marshal tenant: %w", err)
	}
	err = e.cache.Set(e.ctx, key, data, _Tenant_TTL)
	if err != nil {
		return fmt.Errorf("failed to cache tenant: %w", err)
	}
	return nil
}

func (e *Engine) GetCachedTenant(tenantID int32) (*model.Tenant, error) {
	key := tenantKey(tenantID)
	data, err := e.cache.Get(e.ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	var tenant model.Tenant
	err = json.Unmarshal([]byte(data), &tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tenant: %w", err)
	}
	return &tenant, nil
}

func (e *Engine) CacheProduct(productID int32, product *model.Product) error {
	key := productKey(productID)
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}
	err = e.cache.Set(e.ctx, key, data, _PRODUCT_TTL)
	if err != nil {
		return fmt.Errorf("failed to cache product: %w", err)
	}
	return nil
}

func (e *Engine) GetCachedProduct(productID int32) (*model.Product, error) {
	key := productKey(productID)
	data, err := e.cache.Get(e.ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	var product model.Product
	err = json.Unmarshal([]byte(data), &product)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}
	return &product, nil
}

func (e *Engine) CacheBranch(branchID int32, branch *model.Branch) error {
	key := branchKey(branchID)
	data, err := json.Marshal(branch)
	if err != nil {
		return fmt.Errorf("failed to marshal branch: %w", err)
	}
	err = e.cache.Set(e.ctx, key, data, _BRANCH_TTL)
	if err != nil {
		return fmt.Errorf("failed to cache branch: %w", err)
	}
	return nil
}

func (e *Engine) CacheTotalSalesPerProduct(tenantID, productID int32, totalSaleAmount float64) error {
	key := tenantTotalSalesKey(tenantID, productID)

	_, err := e.cache.IncrByFloat(e.ctx, key, totalSaleAmount)
	if err != nil {
		return fmt.Errorf("failed to cache total sales for tenant %d, product %d: %w", tenantID, productID, err)
	}

	err = e.cache.Expire(e.ctx, key, _SALES_PER_PRODUCT_TTL)
	if err != nil {
		return fmt.Errorf("failed to set expiration for total sales cache for tenant %d, product %d: %w", tenantID, productID, err)
	}
	return nil
}

func (e *Engine) GetCachedTotalSalesPerProduct(tenantID, productID int32) (float64, error) {
	key := tenantTotalSalesKey(tenantID, productID)
	return e.cache.GetFloat(e.ctx, key)
}

func (e *Engine) CacheTopSellingProducts(tenantID int32, products []model.ProductSales) error {
	key := topSellingProductsKey(tenantID)
	data, err := json.Marshal(products)
	if err != nil {
		return fmt.Errorf("failed to marshal top selling products: %s", err)
	}
	err = e.cache.Set(e.ctx, key, data, _TOP_SELLING_PRODUCTS_TTL)
	if err != nil {
		return fmt.Errorf("failed to cache top selling products: %s", err)
	}
	return nil
}

func (e *Engine) GetCachedTopSellingProducts(tenantID int32) ([]model.ProductSales, error) {
	key := topSellingProductsKey(tenantID)
	data, err := e.cache.Get(e.ctx, key)
	if err != nil {
		return nil, err
	}
	var products []model.ProductSales
	err = json.Unmarshal([]byte(data), &products)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal top selling products: %s", err)
	}
	return products, nil
}

func (e *Engine) ExpireTopSellingProducts(tenantID int32) error {
	key := topSellingProductsKey(tenantID)
	return e.cache.Expire(e.ctx, key, 0)
}
