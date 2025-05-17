package config

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type FileWriter struct {
	Enabled    bool   `envconfig:"ENABLED" default:"false" json:"enabled"`
	Filename   string `envconfig:"FILE_NAME" default:"logs/app.log" json:"fileName"`
	MaxSizeMB  int    `envconfig:"MAX_SIZE_MB" default:"1" json:"maxSizeMB"`
	MaxBackups int    `envconfig:"MAX_BACKUPS" default:"5" json:"maxBackups"`
	MaxAgeDays int    `envconfig:"MAX_AGE_DAYS" default:"5" json:"maxAgeDays"`
	Compress   bool   `envconfig:"COMPRESS" default:"true" json:"compress"`
}

type ConsoleWriter struct {
	Enabled bool `envconfig:"ENABLED" default:"false" json:"enabled"`
}

type DefaultWriter struct {
	Enabled bool `envconfig:"ENABLED" default:"true" json:"enabled"`
	Output  io.Writer
}

type Log struct {
	FileWriter    `envconfig:"FILE_WRITER" json:"fileWriter"`
	ConsoleWriter `envconfig:"CONSOLE_WRITER" json:"consoleWriter"`
	DefaultWriter `envconfig:"DEFAULT_WRITER" json:"defaultWriter"`
	IconEnabled   bool   `envconfig:"ICON_ENABLED" default:"true" json:"iconEnabled"`
	Level         string `envconfig:"LEVEL" default:"debug" json:"level"`
}

type Config struct {
	Log         `envconfig:"LOG" json:"log"`
	Application struct {
		Env  string `envconfig:"ENV"  default:"production" json:"env"`
		Name string `envconfig:"NAME" default:"MyApp"      json:"name"`
	} `envconfig:"APP" json:"application"`
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		var cfg Config
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatalf("❌ error loading config, %v", err)
		}
		instance = &cfg
		if instance.DefaultWriter.Output == nil {
			instance.DefaultWriter.Output = os.Stdout
		}
	})
	return instance
}
