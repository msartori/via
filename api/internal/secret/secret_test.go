package secret

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"via/internal/log"
	mock_log "via/internal/log/mock"

	"github.com/stretchr/testify/mock"
)

func TestReadSecret_Success(t *testing.T) {
	// Temp secret file creation
	tmpFile, err := os.CreateTemp("", "secret-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "  supersecret\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Execute
	secret := ReadSecret(tmpFile.Name())

	// Eval
	expected := strings.TrimSpace(content)
	if secret != expected {
		t.Errorf("Expected %q, got %q", expected, secret)
	}
}

func TestReadSecret_Error(t *testing.T) {
	if os.Getenv("TEST_SECRET_FATAL") == "1" {
		mockLog := &mock_log.MockLogger{}
		expectedPath := "nonexistent/file/secret.txt"

		mockLog.On("Fatal", mock.Anything, mock.MatchedBy(func(err error) bool {
			pathErr, ok := err.(*os.PathError)
			return ok && strings.Contains(pathErr.Error(), "nonexistent/file/secret.txt")
		}),
			[]any{"msg", "unable to get secret",
				"secret", expectedPath},
		).Run(func(args mock.Arguments) {
			os.Exit(1)
		}).Once()

		log.Set(mockLog)

		ReadSecret(expectedPath)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReadSecret_Error")
	cmd.Env = append(os.Environ(), "TEST_SECRET_FATAL=1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
		t.Fatalf("Expected exit code 1, got: %v, output: %s", err, out.String())
	}
}
