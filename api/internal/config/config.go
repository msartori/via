package config

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v10"
)

type FileWriter struct {
	Enabled    bool   `env:"ENABLED"     envDefault:"true"  json:"enabled"`
	Filename   string `env:"FILE_NAME"   envDefault:"logs/app.log" json:"fileName"`
	MaxSizeMB  int    `env:"MAX_SIZE_MB" envDefault:"1"     json:"maxSizeMB"`
	MaxBackups int    `env:"MAX_BACKUPS" envDefault:"5"     json:"maxBackups"`
	MaxAgeDays int    `env:"MAX_AGE_DAYS" envDefault:"5"    json:"maxAgeDays"`
	Compress   bool   `env:"COMPRESS"    envDefault:"true"  json:"compress"`
}

type ConsoleWriter struct {
	Enabled bool `env:"ENABLED" envDefault:"false" json:"enabled"`
}

type DefaultWriter struct {
	Enabled bool      `env:"ENABLED" envDefault:"true" json:"enabled"`
	Output  io.Writer `json:"-"`
}

type Log struct {
	FileWriter    FileWriter    `envPrefix:"FILE_WRITER_"    json:"fileWriter"`
	ConsoleWriter ConsoleWriter `envPrefix:"CONSOLE_WRITER_" json:"consoleWriter"`
	DefaultWriter DefaultWriter `envPrefix:"DEFAULT_WRITER_" json:"defaultWriter"`
	IconEnabled   bool          `env:"ICON_ENABLED"          envDefault:"true"  json:"iconEnabled"`
	Level         string        `env:"LEVEL"                 envDefault:"debug" json:"level"`
}

type Application struct {
	Env  string `env:"ENV"  envDefault:"production" json:"env"`
	Name string `env:"NAME" envDefault:"via"      json:"name"`
	Port int    `env:"PORT" envDefault:"8080"       json:"port"`
}

type Database struct {
	PasswordFile string `env:"PASSWORD_FILE" envDefault:"" json:"passwordFile"`
}

type CORS struct {
	Enabled bool   `env:"ENABLED" envDefault:"true" json:"enabled"`
	Origins string `env:"ORIGINS" envDefault:"*" json:"origins"`
	Methods string `env:"METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS" json:"methods"`
	Headers string `env:"HEADERS" envDefault:"Content-Type,Authorization,bypass-tunnel-reminder" json:"headers"`
}

type Config struct {
	Log         Log         `envPrefix:"LOG_" json:"log"`
	Application Application `envPrefix:"APP_" json:"application"`
	Database    Database    `envPrefix:"DB_" json:"db"`
	CORS        CORS        `envPrefix:"CORS_" json:"cors"`
}

var (
	instance *Config
	once     sync.Once
)

// Get returns a singleton config loaded from environment variables
func Get() *Config {
	once.Do(func() {
		var cfg Config
		opts := env.Options{
			Prefix:          "",
			TagName:         "env",
			RequiredIfNoDef: false,
		}
		if err := env.ParseWithOptions(&cfg, opts); err != nil {
			log.Fatalf("❌ Error loading config: %v", err)
		}

		if cfg.Log.DefaultWriter.Output == nil {
			cfg.Log.DefaultWriter.Output = os.Stdout
		}
		instance = &cfg
	})
	return instance
}
