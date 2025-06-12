package app_log

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"via/internal/log"
)

var logTrace = "\"key\":\"trace\""
var logDebug = "\"key\":\"debug\""
var logInfo = "\"key\":\"info\""
var logWarn = "\"key\":\"warn\""
var logError = "\"key\":\"error\""
var logFatal = "\"key\":\"fatal\""

func getBaseLogConfig() LogCfg {
	return LogCfg{
		IconEnabled: false,
		DefaultWriter: DefaultWriterCfg{
			Enabled: true,
		},
		FileWriter: FileWriterCfg{
			Enabled: false,
		},
	}
}

func logAll(logger log.Logger, logLvl string) {
	key := "key"
	ctx := context.Background()
	if logLvl == lvlTrace {
		logger.Trace(ctx, key, lvlTrace)
		logger.Debug(ctx, key, lvlDebug)
		logger.Info(ctx, key, lvlInfo)
		logger.Warn(ctx, key, lvlWarn)
		logger.Error(ctx, errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlDebug {
		logger.Debug(ctx, key, lvlDebug)
		logger.Info(ctx, key, lvlInfo)
		logger.Warn(ctx, key, lvlWarn)
		logger.Error(ctx, errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlInfo {
		logger.Info(ctx, key, lvlInfo)
		logger.Warn(ctx, key, lvlWarn)
		logger.Error(ctx, errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlWarn {
		logger.Warn(ctx, key, lvlWarn)
		logger.Error(ctx, errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlError {
		logger.Error(ctx, errors.New(lvlError), key, lvlError)
	}
}

func evalLogNotFound(t *testing.T, result string, logToFind string) {
	if !strings.Contains(result, logToFind) {
		t.Errorf("Log with expected level not found, expected:%s, got:%s", logToFind, result)
	}
}

func evalLogFound(t *testing.T, result string, logToFind string) {
	if strings.Contains(result, logToFind) {
		t.Errorf("Log with Not expected level found, not expected:%s, got:%s", logToFind, result)
	}
}

func TestFatal(t *testing.T) {
	key := "key"
	if os.Getenv("TEST_FATAL") == "1" {
		cfg := getBaseLogConfig()
		cfg.Level = lvlError
		cfg.DefaultWriter.Output = os.Stdout
		log.Set(New(cfg)) // Initialize the logger with the config
		logger := log.Get()
		logAll(logger, lvlTrace)
		logger.Fatal(context.Background(), errors.New("fatal"), key, lvlFatal)
		return
	}
	var logOutput bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	cmd.Stdout = &logOutput
	cmd.Stderr = &logOutput
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		result := (logOutput.String())
		evalLogFound(t, result, logTrace)
		evalLogFound(t, result, logDebug)
		evalLogFound(t, result, logInfo)
		evalLogFound(t, result, logWarn)
		evalLogNotFound(t, result, logError)
		evalLogNotFound(t, result, logFatal)
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}

func TestFileWriter(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testlog-*.log")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	cfg := getBaseLogConfig()
	cfg.Level = lvlInfo
	cfg.FileWriter.Enabled = true
	cfg.FileWriter.Filename = tmpFile.Name()
	cfg.DefaultWriter.Enabled = false

	log.Set(New(cfg))
	logger := log.Get()
	logger.Info(context.Background(), "key", "filewriter")

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading log file: %v", err)
	}

	if !strings.Contains(string(content), "filewriter") {
		t.Errorf("Expected log not found in file, got: %s", string(content))
	}
}

func TestConsoleWriter(t *testing.T) {
	// Backup original stdout
	originalStdout := os.Stdout

	// Create a pipe to capture stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Redirect stdout
	os.Stdout = writePipe

	// Setup logger with console writer
	cfg := getBaseLogConfig()
	cfg.Level = lvlInfo
	cfg.ConsoleWriter.Enabled = true
	cfg.DefaultWriter.Enabled = false
	log.Set(New(cfg))
	logger := log.Get()

	// Log something
	logger.Info(context.Background(), "key", "consolewriter")

	// Close writer to finish the pipe and restore stdout
	writePipe.Close()
	os.Stdout = originalStdout // Restore stdout

	// Read captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(readPipe)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "consolewriter") {
		t.Errorf("Expected output not found in console log, got: %s", output)
	}
}

func TestUnknownLogLevel(t *testing.T) {
	if os.Getenv("TEST_UNKNOWN_LEVEL") == "1" {
		cfg := getBaseLogConfig()
		cfg.Level = "nonexistent"
		New(cfg)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestUnknownLogLevel")
	cmd.Env = append(os.Environ(), "TEST_UNKNOWN_LEVEL=1")
	err := cmd.Run()

	if err == nil {
		t.Fatalf("Expected fatal error due to unknown log level")
	}
}

func TestTraceLevelEnabled(t *testing.T) {
	var buf bytes.Buffer
	cfg := getBaseLogConfig()
	cfg.Level = lvlTrace
	cfg.DefaultWriter.Output = &buf
	log.Set(New(cfg))
	logger := log.Get()

	logger.Trace(context.Background(), "key", "tracelevel")

	if !strings.Contains(buf.String(), "tracelevel") {
		t.Errorf("Expected 'tracelevel' in output, got: %s", buf.String())
	}
}

func TestLoggerWithNoWriters(t *testing.T) {
	cfg := getBaseLogConfig()
	cfg.Level = lvlInfo
	cfg.DefaultWriter.Enabled = false
	cfg.ConsoleWriter.Enabled = false
	cfg.FileWriter.Enabled = false

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger panicked with no writers configured: %v", r)
		}
	}()

	_ = New(cfg) // Just ensure it doesn't panic
}

func TestLoggerWithIconsEnabled(t *testing.T) {
	var buf bytes.Buffer
	cfg := getBaseLogConfig()
	cfg.IconEnabled = true
	cfg.Level = lvlInfo
	cfg.DefaultWriter.Output = &buf

	log.Set(New(cfg))
	logger := log.Get()
	logger.Info(context.Background(), "key", "with-icon")

	output := buf.String()
	if !strings.Contains(output, "with-icon") {
		t.Errorf("Expected 'with-icon' log message in output, got: %s", output)
	}
	if !strings.Contains(output, "\"lvl\":\"ℹ️\"") {
		t.Errorf("Expected info icon in log output: %s", output)
	}
}

func TestAddFields_InvalidKeyAndOddArgs(t *testing.T) {
	var buf bytes.Buffer

	cfg := getBaseLogConfig()
	cfg.Level = lvlWarn
	cfg.DefaultWriter.Output = &buf

	log.Set(New(cfg))
	log.Get().Info(context.Background(), 123, "valueOfNonStringIgnoredKey", "keyWithoutValue")

	output := buf.String()

	// Validate keyWithoutValue is ignored due to missing value
	if strings.Contains(output, "\"keyWithoutValue\":") {
		t.Errorf("Expected 'keyWithoutValue' to be ignored due to missing value, got: %s", output)
	}

	// Validate valueOfNonStringIgnoredKey is skipped due to non-string key
	if strings.Contains(output, "valueOfNonStringIgnoredKey") {
		t.Errorf("Expected valueOfNonStringIgnoredKey to be skipped due to non-string key, got: %s", output)
	}

	// Validate warning about odd number of arguments
	if !strings.Contains(output, "Logger received an odd number of arguments. Last argument ignored.") {
		t.Errorf("Expected warning about odd number of arguments, got: %s", output)
	}
}

func TestLogFieldsHelpers(t *testing.T) {
	var buf bytes.Buffer

	cfg := getBaseLogConfig()
	cfg.Level = lvlWarn
	cfg.DefaultWriter.Output = &buf

	log.Set(New(cfg))
	logger := log.Get()

	// Casos para WithLogFields
	ctx := context.Background()
	ctx = logger.WithLogFields(ctx, "k1", "v1", 123, "ignored", "orphan") // 123 no es string, "orphan" es impar

	// Verifica que "k1" fue añadido
	fields, ok := ctx.Value(logContextKey).(LogFields)
	if !ok || fields == nil {
		t.Fatal("Expected LogFields in context")
	}
	if val, exists := fields["k1"]; !exists || val != "v1" {
		t.Errorf("Expected key k1 with value v1, got: %v", fields)
	}

	// Verifica que 123 no fue añadido como clave y "orphan" fue ignorado
	if _, exists := fields["ignored"]; exists {
		t.Errorf("Expected key 'ignored' to be skipped, got: %v", fields)
	}
	if _, exists := fields["orphan"]; exists {
		t.Errorf("Expected 'orphan' to be ignored due to odd count, got: %v", fields)
	}
	if !strings.Contains(buf.String(), "Logger received an odd number of arguments") {
		t.Errorf("Expected warning about odd number of arguments")
	}

	// Casos para WithLogFieldsInRequest
	req, _ := http.NewRequest("GET", "/", nil)
	req = logger.WithLogFieldsInRequest(req, "rk", "rv")

	rctx := req.Context().Value(logContextKey).(LogFields)
	if val, ok := rctx["rk"]; !ok || val != "rv" {
		t.Errorf("Expected request context field rk=rv, got: %v", rctx)
	}
	ctx = logger.WithLogFields(nil, "key", "value") // Test with nil context, should not panic
	if ctx != nil {
		t.Errorf("Expected nil context to return nil, got: %v", ctx)
	}

	logger.Warn(nil, "keyCtxNil", "value") // Test with nil context, should not panic
	output := buf.String()
	if !strings.Contains(output, "keyCtxNil") {
		t.Errorf("Expected 'keyCtxNil' in output, got: %s", output)
	}

	ctx = context.Background()
	ctx = logger.WithLogFields(ctx, "addedToCtx", "value") // Test with valid context and pairs
	logger.Warn(ctx, "key", "value")
	output = buf.String()
	if !strings.Contains(output, "addedToCtx") {
		t.Errorf("Expected 'addedToCtx' in output, got: %s", output)
	}

	/*

		// Casos para addContextFields
		merged := appLogger.addContextFields(req.Context(), "x", "y")

		expected := map[string]any{
			"rk": "rv",
			"x":  "y",
		}

		if len(merged)%2 != 0 {
			t.Errorf("Expected even number of merged fields, got: %v", merged)
		}

		m := make(map[string]any)
		for i := 0; i < len(merged); i += 2 {
			k, ok := merged[i].(string)
			if !ok {
				t.Errorf("Unexpected non-string key: %v", merged[i])
				continue
			}
			m[k] = merged[i+1]
		}

		for k, v := range expected {
			if m[k] != v {
				t.Errorf("Expected key %s with value %v, got %v", k, v, m[k])
			}
		}
	*/
}
