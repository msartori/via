package main

import (
	"context"
	"via/internal/app"
	ent_client "via/internal/client/ent"
	"via/internal/config"
	db_pool "via/internal/db/pool"
	"via/internal/log"
	app_log "via/internal/log/app"
)

func main() {
	cfg := config.Get()
	log.Set(app_log.New(cfg.Log))

	logger := log.Get()
	logger.Debug(context.Background(), "config", cfg)

	dbPool, err := db_pool.New(cfg.Database)
	if err != nil {
		logger.Fatal(context.Background(), err, "msg", "Error al conectar")
	}
	entClient := ent_client.New(dbPool)
	// Run the auto migration tool.
	if err := entClient.Schema.Create(context.Background()); err != nil {
		logger.Fatal(context.Background(), err, "msg", "failed creating schema resources")
	}
	defer entClient.Close()
	defer dbPool.Close()

	app.StartServer(cfg)
}
