package app_log

import (
	"io"
	"os"
	"time"
	"via/internal/config"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileWriter struct {
	config config.FileWriter
}

func (fw FileWriter) Writer() io.Writer {
	return &lumberjack.Logger{
		Filename:   fw.config.Filename,
		MaxSize:    fw.config.MaxSizeMB,
		MaxBackups: fw.config.MaxBackups,
		MaxAge:     fw.config.MaxAgeDays,
		Compress:   fw.config.Compress,
	}
}

type ConsoleWriter struct {
}

func (cw ConsoleWriter) Writer() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
}

type DefaultWriter struct {
	config config.DefaultWriter
}

func (dw DefaultWriter) Writer() io.Writer {
	return dw.config.Output
}
