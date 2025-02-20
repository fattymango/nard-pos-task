package main

import (
	"fmt"
	"multitenant/app"
	"multitenant/pkg/config"
	"multitenant/pkg/db"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %s", err))
	}

	db, err := db.NewDB(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to connect to db: %s", err))
	}

	app.Migrate(cfg, db)
}
