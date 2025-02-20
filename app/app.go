package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"multitenant/handler"
	"multitenant/pkg/cache"
	"multitenant/pkg/config"
	"multitenant/pkg/db"
	"multitenant/pkg/logger"
	"multitenant/pkg/metrics"
	"multitenant/pkg/rabbitmq"
)

func Start() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %s", err))
	}

	log, err := logger.NewLogger(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %s", err))
	}

	// DB connection
	log.Info("Creating db connection...")

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to create db connection: %s", err)
	}
	log.Info("DB connection successful")

	// Auto migrate
	log.Info("Auto migrating...")
	err = Migrate(cfg, db)
	if err != nil {
		log.Fatalf("failed to auto migrate: %s", err)
	}
	log.Info("Auto migration successful")

	// Seed, for demo purposes, shouild be independent of the program
	log.Info("Seeding...")
	err = Seed(db)
	if err != nil {
		log.Errorf("failed to seed: %s", err)
		log.Info("Continuing...")
	}
	log.Info("Seed successful")

	// Cache connection
	log.Info("Creating cache connection...")
	cache, err := cache.NewCache(cfg)
	if err != nil {
		log.Fatalf("failed to create cache: %s", err)
	}
	log.Info("Cache connection successful")

	// RabbitMQ connection
	log.Info("Creating RabbitMQ connection...")
	amqp, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ connection: %s", err)
	}
	log.Info("RabbitMQ connection successful")

	// Metrics
	metrics := metrics.NewMetrics(ctx, cfg)
	metrics.StartMonitoring()

	log.Info("Creating Multi Tenant handler...")
	server, err := handler.NewMultiTanentServer(ctx, cfg, log, db, cache, amqp, metrics)
	if err != nil {
		log.Fatalf("failed to create Multi Tenant handler: %s", err)
	}
	log.Info("Multi Tenant handler created")

	log.Info("Starting Multi Tenant server...")
	err = server.Start()
	if err != nil {
		log.Fatalf("failed to start Multi Tenant server: %s", err)
	}

	if err := recover(); err != nil {
		log.Fatalf("some panic...: %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Infof("signal.Notify CTRL+C: %v", v)
	}

	log.Info("Shutting down Multi Tenant server...")

	cancel()
	server.Stop()

}
