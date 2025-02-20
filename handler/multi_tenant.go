package handler

import (
	"context"
	"fmt"
	"multitenant/internal/engine"
	"multitenant/pkg/cache"
	"multitenant/pkg/config"
	"multitenant/pkg/db"
	"multitenant/pkg/logger"
	"multitenant/pkg/metrics"
	"multitenant/pkg/rabbitmq"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MultiTanentServer struct {
	ctx       context.Context
	config    *config.Config
	logger    *logger.Logger
	db        *db.DB
	cache     *cache.Cache
	validator *validator.Validate
	metrics   *metrics.Metrics
	// Server // App Interfaces
	amqp *rabbitmq.RabbitMQ
	app  *fiber.App
	pb   *MultiTenantRPCServer

	engine *engine.Engine
}

func NewMultiTanentServer(ctx context.Context, cfg *config.Config, logger *logger.Logger, db *db.DB, cache *cache.Cache, amqp *rabbitmq.RabbitMQ, m *metrics.Metrics) (*MultiTanentServer, error) {

	e := engine.NewEngine(ctx, cfg, logger, db, cache, m)

	pb := NewMultiTenantRPCServer(ctx, cfg, logger, e)

	return &MultiTanentServer{
		ctx:       ctx,
		config:    cfg,
		logger:    logger,
		db:        db,
		cache:     cache,
		validator: validator.New(),
		metrics:   m,
		engine:    e,

		amqp: amqp,
		app:  fiber.New(),
		pb:   pb,
	}, nil
}

func (s *MultiTanentServer) Start() error {
	s.RegisterRoutes()

	go func() {
		if err := s.ConsumeTransactions(); err != nil {
			s.logger.Fatalf("failed to consume transactions: %s", err)
		}
	}()

	go func() {
		if err := s.pb.Start(); err != nil {
			s.logger.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	go func() {
		err := s.engine.Start()
		if err != nil {
			s.logger.Fatalf("failed to start engine: %s", err) // fatal here, as engine is critical
		}

		if err := s.app.Listen(fmt.Sprintf(":%s", s.config.Server.Port)); err != nil {
			s.logger.Fatalf("failed to start server: %s", err) // fatal here, as server is critical
		}
	}()
	return nil

}

func (s *MultiTanentServer) Stop() {
	err := s.amqp.Close()
	if err != nil {
		s.logger.Errorf("failed to close RabbitMQ connection: %s", err)
	}

	err = s.app.Shutdown()
	if err != nil {
		s.logger.Errorf("failed to shutdown server: %s", err)
	}

	s.pb.Stop()

	s.engine.Stop()

	s.db.Close()
}
