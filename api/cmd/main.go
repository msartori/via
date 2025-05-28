package main

import (
	"fmt"
	"net/http"
	"via/internal/config"
	handler_cors "via/internal/handler/cors"
	handler_guide "via/internal/handler/guide"
	log "via/internal/util/logger"
)

func main() {
	logger := log.Get()

	http.HandleFunc("/guide", handler_guide.GetGuide())
	handler := handler_cors.WithCORS(http.DefaultServeMux)
	logger.Info("config", config.Get())
	logger.Info("msg", "Servidor iniciado en http://localhost:8080", "")
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Get().Application.Port), handler))
}
