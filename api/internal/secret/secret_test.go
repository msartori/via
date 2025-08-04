package secret

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	mock_secret "via/internal/secret/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetAndGet(t *testing.T) {
	mock := new(mock_secret.MockSecret)
	mock.On("Read", "any-path").Return("mocked-secret")
	Set(mock)

	got := Get()
	assert.Equal(t, "mocked-secret", got.Read("any-path"))
	mock.AssertExpectations(t)
}

func TestFileSecretReader_Read_Success(t *testing.T) {
	tmpDir := t.TempDir()
	secretFile := filepath.Join(tmpDir, "secret.txt")
	err := os.WriteFile(secretFile, []byte("  super-secret\n"), 0600)
	assert.NoError(t, err)

	reader := &FileSecretReader{}
	secretValue := reader.Read(secretFile)
	assert.Equal(t, "super-secret", secretValue)
}

func TestFileSecretReader_Read_Error_TriggersFatal(t *testing.T) {
	// Mock logger and replace global log
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger)

	ctx := context.Background()
	path := "/nonexistent/secret"
	mockLogger.On("Fatal", ctx, mock.Anything, []any{"msg", "unable to get secret", "secret", path}).Return()

	reader := &FileSecretReader{}
	reader.Read(path) // will trigger logger.Fatal

	mockLogger.AssertCalled(t, "Fatal", ctx, mock.Anything, []any{"msg", "unable to get secret", "secret", path})
}
