package main

import (
	"context"
	"net/http"
	"via/internal/auth"
	ent_client "via/internal/client/ent"
	"via/internal/config"
	db_pool "via/internal/db/pool"
	"via/internal/ds"
	redis_ds "via/internal/ds/redis"
	jwt_key "via/internal/jwt"
	"via/internal/log"
	app_log "via/internal/log/app"
	guide_provider "via/internal/provider/guide"
	guide_ent_provider "via/internal/provider/guide/ent"
	operator_provider "via/internal/provider/operator"
	operator_ent_provider "via/internal/provider/operator/ent"
	via_guide_provider "via/internal/provider/via/guide"
	via_guide_web_provider "via/internal/provider/via/guide/web"
	"via/internal/pubsub"
	redis_pubsub "via/internal/pubsub/redis"
	"via/internal/router"
	"via/internal/secret"
	"via/internal/server"
)

func main() {
	// Load the configuration
	cfg := config.Get()

	//log initialization
	log.Set(app_log.New(cfg.Log))
	logger := log.Get()
	logger.Debug(context.Background(), "config", cfg)

	secret.Set(new(secret.FileSecretReader))

	ds.Set(redis_ds.New(cfg.DS))

	pubsub.Set(redis_pubsub.New(cfg.PubSub))

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

	// Set up dependencies
	via_guide_provider.Set(via_guide_web_provider.New(cfg.GuideWebClient, via_guide_web_provider.HistoricalQueryResponseParser{}))
	guide_provider.Set(guide_ent_provider.New())
	operator_provider.Set(operator_ent_provider.New())

	servers := []*http.Server{}
	//start servers
	if cfg.RestServer.Enabled {
		servers = append(servers, server.Create(cfg.RestServer, router.NewRest(cfg)))
	}
	if cfg.SSEServer.Enabled {
		servers = append(servers, server.Create(cfg.SSEServer, router.NewSSE(cfg)))
	}
	server.Start(servers)
}
