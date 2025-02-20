package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *MultiTanentServer) Middleware_Logger(ctx *fiber.Ctx) error {
	start := time.Now()
	err := ctx.Next()
	duration := time.Since(start)

	if err != nil {
		s.logger.Errorf("Error: %s", err.Error())
	}

	s.logger.Infof("%d | %s | %s | %v", ctx.Response().StatusCode(), ctx.Method(), ctx.Path(), duration)
	return err
}

func (s *MultiTanentServer) Middleware_MetricsDuration(ctx *fiber.Ctx) error {
	now := time.Now()
	err := ctx.Next()
	s.metrics.SetAPIRequestDuration(ctx.Path(), ctx.Method(), time.Since(now).Seconds())
	return err
}
