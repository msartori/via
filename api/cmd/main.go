package main

import (
	"context"
	"via/internal/app"
	"via/internal/auth"
	ent_client "via/internal/client/ent"
	"via/internal/config"
	db_pool "via/internal/db/pool"
	"via/internal/ds"
	redis_ds "via/internal/ds/redis"
	jwt_key "via/internal/jwt"
	"via/internal/log"
	app_log "via/internal/log/app"
)

func main() {
	// Load the configuration
	cfg := config.Get()

	//log initialization
	log.Set(app_log.New(cfg.Log))
	logger := log.Get()
	logger.Debug(context.Background(), "config", cfg)

	ds.Set(redis_ds.New(cfg.DS))
	auth.Set(auth.New(cfg.OAuth, ds.Get()))

	// JWT key initialization
	if err := jwt_key.Init(cfg.JWT); err != nil {
		logger.Fatal(context.Background(), err, "msg", "Error in jwt key initialization")
	}

	// Initialize the database connection pool
	dbPool, err := db_pool.New(cfg.Database)
	if err != nil {
		logger.Fatal(context.Background(), err, "msg", "Error in db connection pool initialization")
	}
	// Initialize the Ent client with the database connection pool
	entClient := ent_client.New(dbPool)
	defer entClient.Close()
	defer dbPool.Close()

	app.StartServer(cfg)
}
