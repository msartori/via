package config

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_NAME", "myapp")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("APP_REQUEST_TIMEOUT", "45")
	t.Setenv("LOG_DEFAULTWRITER_OUTPUT", "stdout") // dummy env

	reset()

	cfg := Get()
	assert.Equal(t, "test", cfg.Application.Env)
	assert.Equal(t, "myapp", cfg.Application.Name)
	assert.Equal(t, 9090, cfg.Application.Port)
	assert.Equal(t, 45, cfg.Application.RequestTimeout)

	// Check singleton behavior
	cfg2 := Get()
	assert.Equal(t, cfg, cfg2)
}

func TestConfigFallbackOutput(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_NAME", "myapp")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("APP_REQUEST_TIMEOUT", "45")

	reset()

	cfg := Get()
	assert.NotNil(t, cfg.Log.DefaultWriter.Output)
}

func TestGet_ErrorCase(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		reset()
		// Set invalid int value for APP_PORT
		os.Setenv("APP_PORT", "invalid_port")
		defer os.Unsetenv("APP_PORT")
		_ = Get() // Should call log.Fatalf (which triggers os.Exit)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestGet_ErrorCase")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected subprocess to exit with error")
	}

	if !strings.Contains(string(output), "❌ Error loading config") {
		t.Errorf("expected fatal error log, got: %s", string(output))
	}
}
