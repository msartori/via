package util_logger

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"via/internal/config"
)

var logTrace = "\"key\":\"trace\""
var logDebug = "\"key\":\"debug\""
var logInfo = "\"key\":\"info\""
var logWarn = "\"key\":\"warn\""
var logError = "\"key\":\"error\""
var logFatal = "\"key\":\"fatal\""

func getBaseLogConfig() config.Log {
	return config.Log{
		IconEnabled: false,
		DefaultWriter: config.DefaultWriter{
			Enabled: true,
		},
		FileWriter: config.FileWriter{
			Enabled: false,
		},
	}
}

func logAll(logger *Log, logLvl string) {
	key := "key"
	if logLvl == lvlTrace {
		logger.Trace(key, lvlTrace)
		logger.Debug(key, lvlDebug)
		logger.Info(key, lvlInfo)
		logger.Warn(key, lvlWarn)
		logger.Error(errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlDebug {
		logger.Debug(key, lvlDebug)
		logger.Info(key, lvlInfo)
		logger.Warn(key, lvlWarn)
		logger.Error(errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlInfo {
		logger.Info(key, lvlInfo)
		logger.Warn(key, lvlWarn)
		logger.Error(errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlWarn {
		logger.Warn(key, lvlWarn)
		logger.Error(errors.New(lvlError), key, lvlError)
	}
	if logLvl == lvlError {
		logger.Error(errors.New(lvlError), key, lvlError)
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
		Init(cfg)
		logger := Get()
		logAll(logger, lvlTrace)
		logger.Fatal(errors.New("fatal"), key, lvlFatal)
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

/*
func TestLoggerLevel(t *testing.T) {
	var logOutput bytes.Buffer
	cfg := getBaseLogConfig()

	Init(cfg)
	logger := Get()
	logger.Trace("key", "trace")
	logger.Debug("key", "debug")
	logger.Info("key", "info")
	logger.Warn("key", "warn")
	logger.Error(errors.New("error"), "key", "error")
	output := logOutput.String()
	if !strings.Contains(output, "fatal") || strings.Contains(output, "trace") {
		t.Errorf("Expected log output to contain provided key-value pairs, got: %s", output)
	}
}

func TestLogger_Info(t *testing.T) {
	var logOutput bytes.Buffer
	cfg := config.Log{
		Level:       "info",
		IconEnabled: true,
		DefaultWriter: config.DefaultWriter{
			Enabled: true,
			Output:  &logOutput,
		},
		FileWriter: config.FileWriter{
			Enabled: false,
		},
	}
	Init(cfg)
	logger := Get()
	logger.Info("key1", "value1", "key2", "value")
	output := logOutput.String()
	if !strings.Contains(output, "value1") || !strings.Contains(output, "value2") {
		t.Errorf("Expected log output to contain provided key-value pairs, got: %s", output)
	}
}
*/
