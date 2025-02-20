package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *MultiTanentServer) GetTotalSalesPerProduct(ctx *fiber.Ctx) error {
	s.logger.Debug("GetTotalSalesPerProduct")
	tenantID, err := ctx.ParamsInt("tenantID")
	if err != nil {
		return NewBadRequestResponse(ctx, "invalid tenantID")
	}

	productID, err := ctx.ParamsInt("productID")
	if err != nil {
		return NewBadRequestResponse(ctx, "invalid productID")
	}

	ok, err := s.engine.TanentExists(int32(tenantID))
	if err != nil {
		return NewInternalServerErrorResponse(ctx, fmt.Sprintf("failed to check tenant existence: %s", err))
	}

	if !ok {
		return NewNotFoundResponse(ctx, "tenant not found")
	}

	ok, err = s.engine.ProductExists(int32(productID))
	if err != nil {
		return NewInternalServerErrorResponse(ctx, fmt.Sprintf("failed to check product existence: %s", err))
	}

	if !ok {
		return NewNotFoundResponse(ctx, "product not found")
	}

	totalSales, err := s.engine.GetTotalSalesPerProduct(int32(tenantID), int32(productID))
	if err != nil {
		return NewInternalServerErrorResponse(ctx, err.Error())
	}

	return NewSuccessResponse(ctx, totalSales)
}

func (s *MultiTanentServer) GetTopSellingProducts(ctx *fiber.Ctx) error {
	s.logger.Debug("GetTopSellingProducts")
	tenantID, err := ctx.ParamsInt("tenantID")
	if err != nil {
		return NewBadRequestResponse(ctx, "invalid tenantID")
	}

	ok, err := s.engine.TanentExists(int32(tenantID))
	if err != nil {
		return NewInternalServerErrorResponse(ctx, fmt.Sprintf("failed to check tenant existence: %s", err))
	}

	if !ok {
		return NewNotFoundResponse(ctx, "tenant not found")
	}

	products, err := s.engine.GetTopSellingProducts(int32(tenantID))
	if err != nil {
		return NewInternalServerErrorResponse(ctx, fmt.Sprintf("failed to get top selling products: %s", err))
	}

	return NewSuccessResponse(ctx, products)

}
