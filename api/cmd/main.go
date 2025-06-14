package main

import (
	"via/internal/app"
	"via/internal/config"
	"via/internal/log"
	app_log "via/internal/log/app"
)

func main() {
	cfg := config.Get()
	log.Set(app_log.New(cfg.Log))
	app.StartServer(cfg)
}
