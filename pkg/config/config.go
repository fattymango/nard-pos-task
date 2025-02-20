package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Logger struct {
	LogFile string `envconfig:"LOG_FILE" validate:"required"`
}

type Server struct {
	Port string `envconfig:"SERVER_PORT" validate:"required,numeric"`
}

type Debug struct {
	Debug bool `envconfig:"DEBUG"`
}

type DB struct {
	Host            string `envconfig:"MYSQL_HOST" validate:"required"`
	Port            string `envconfig:"MYSQL_PORT" validate:"required"`
	User            string `envconfig:"MYSQL_USER" validate:"required"`
	Password        string `envconfig:"MYSQL_PASSWORD" validate:"required"`
	Name            string `envconfig:"MYSQL_DATABASE" validate:"required"`
	SSLMode         string `envconfig:"MYSQL_SSL_MODE" validate:"required"`
	MaxIdleConns    int    `envconfig:"MYSQL_MAX_IDLE_CONNS" default:"2"`
	MaxOpenConns    int    `envconfig:"MYSQL_MAX_OPEN_CONNS" default:"5"`
	MaxConnLifetime int    `envconfig:"MYSQL_MAX_CONN_LIFETIME" default:"10"`
}

type Redis struct {
	Host string `envconfig:"REDIS_HOST" validate:"required"`
	Port string `envconfig:"REDIS_PORT" validate:"required"`
}

type RabbitMQ struct {
	Protocol string `envconfig:"AMQP_PROTOCOL" validate:"required"`
	Host     string `envconfig:"AMQP_HOST" validate:"required"`
	Port     string `envconfig:"AMQP_PORT" validate:"required"`
	Username string `envconfig:"AMQP_USERNAME" validate:"required"`
	Password string `envconfig:"AMQP_PASSWORD" validate:"required"`
}

type GRPC struct {
	Host string `envconfig:"GRPC_HOST" validate:"required"`
	Port string `envconfig:"GRPC_PORT" validate:"required"`
}

type Config struct {
	DB       DB
	Logger   Logger
	Server   Server
	Debug    Debug
	Redis    Redis
	RabbitMQ RabbitMQ
	GRPC     GRPC
}

func NewConfig() (*Config, error) {
	fmt.Println("Loading configuration...")

	debug := os.Getenv("DEBUG")
	fmt.Printf("DEBUG: %s\n", debug)

	cfg := &Config{}

	// Load environment variables into the Config struct using envconfig
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate the config struct
	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Print out the loaded config (for testing purposes)
	log.Printf("Configuration Loaded: %+v\n\n", cfg)
	return cfg, nil
}

func NewTestConfig() (*Config, error) {
	cfg := &Config{
		Server: Server{
			Port: "8080",
		},
	}

	return cfg, nil
}
