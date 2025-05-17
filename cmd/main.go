package main

import (
	"net/http"
	"via/internal/config"
	handler_guide "via/internal/handler/guide"
	log "via/internal/util/logger"
)

func main() {
	logger := log.Get()

	http.HandleFunc("/guide", handler_guide.GetGuide())
	logger.Info("config", config.Get())
	logger.Info("msg", "Servidor iniciado en http://localhost:8080", "")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
