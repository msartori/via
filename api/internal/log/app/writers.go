package app_log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileWriterCfg struct {
	Enabled    bool   `env:"ENABLED"     envDefault:"true"  json:"enabled"`
	Filename   string `env:"FILE_NAME"   envDefault:"logs/app.log" json:"fileName"`
	MaxSizeMB  int    `env:"MAX_SIZE_MB" envDefault:"1"     json:"maxSizeMB"`
	MaxBackups int    `env:"MAX_BACKUPS" envDefault:"5"     json:"maxBackups"`
	MaxAgeDays int    `env:"MAX_AGE_DAYS" envDefault:"5"    json:"maxAgeDays"`
	Compress   bool   `env:"COMPRESS"    envDefault:"true"  json:"compress"`
}

type ConsoleWriterCfg struct {
	Enabled bool `env:"ENABLED" envDefault:"false" json:"enabled"`
}

type DefaultWriterCfg struct {
	Enabled bool      `env:"ENABLED" envDefault:"true" json:"enabled"`
	Output  io.Writer `json:"-"`
}

type FileWriter struct {
	config FileWriterCfg
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
	config DefaultWriterCfg
}

func (dw DefaultWriter) Writer() io.Writer {
	return dw.config.Output
}
