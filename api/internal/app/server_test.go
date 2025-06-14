package app

import (
	"syscall"
	"testing"
	"time"
	"via/internal/config"
	"via/internal/log"
	mock_log "via/internal/log/mock"
)

func TestStartServer_GracefulShutdown(t *testing.T) {
	cfg := config.Config{
		Application: config.Application{
			Port:           8099, // Free test port
			RequestTimeout: 2,
		},
	}

	mockLog := new(mock_log.MockNoOpLogger)
	// logger mock setup
	log.Set(mockLog)

	// Execute server in go routine
	go StartServer(cfg)

	// Waits for server start
	time.Sleep(200 * time.Millisecond)

	// Sends a SIGINT signal to simulate Ctrl+C (shutdown)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	// Waits for shutdown completes
	time.Sleep(500 * time.Millisecond)

}
