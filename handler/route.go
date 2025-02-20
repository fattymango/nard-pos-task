package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func (s *MultiTanentServer) RegisterRoutes() {

	root := s.app.Group("/")
	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	// CORS Middleware
	root.Use(cors.New())

	// Loging Middleware
	root.Use(s.Middleware_Logger)
	// Metrics Middleware
	root.Use(s.Middleware_MetricsDuration)
	// Metrics
	root.Get("/metrics", func(c *fiber.Ctx) error {
		p(c.Context())
		return nil
	})
	api := root.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/transaction", s.CreateTransaction)
	v1.Get("/tenant/:tenantID/product/:productID/sales", s.GetTotalSalesPerProduct)
	v1.Get("/tenant/:tenantID/product/top", s.GetTopSellingProducts)
}
