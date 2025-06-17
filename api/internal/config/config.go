package config

import (
	"log"
	"os"
	"sync"
	http_client "via/internal/client/http"
	app_log "via/internal/log/app"

	"github.com/caarlos0/env/v10"
)

type Bussiness struct {
	ViaBranch       string `env:"VIA_BRANCH"  envDefault:"123" json:"viaBranch"`
	WithdrawStatus  string `env:"WITHDRAW_STATUS"  envDefault:"CRR" json:"withdrawStatus"`
	DeliveredStatus string `env:"DELIVERED_STATUS" envDefault:"ENT" json:"deliveredStatus"`
	PendingStatus   string `env:"PENDING_STATUS"  envDefault:"ASP,CPO,ORI,PTE" json:"pendingStatus"`
}

type Application struct {
	Env            string `env:"ENV"  envDefault:"production" json:"env"`
	Name           string `env:"NAME" envDefault:"via"      json:"name"`
	Port           int    `env:"PORT" envDefault:"8080"       json:"port"`
	RequestTimeout int    `env:"REQUEST_TIMEOUT" envDefault:"30"       json:"requestTimeout"`
}

type Database struct {
	PasswordFile string `env:"PASSWORD_FILE" envDefault:"" json:"passwordFile"`
	User         string `env:"USER" envDefault:"" json:"-"`
	Base         string `env:"BASE" envDefault:"" json:"-"`
	Port         string `env:"PORT" envDefault:"" json:"port"`
	Host         string `env:"HOST" envDefault:"" json:"host"`
}

type CORS struct {
	Enabled bool   `env:"ENABLED" envDefault:"true" json:"enabled"`
	Origins string `env:"ORIGINS" envDefault:"*" json:"origins"`
	Methods string `env:"METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS" json:"methods"`
	Headers string `env:"HEADERS" envDefault:"Content-Type,Authorization,bypass-tunnel-reminder,Accept-Language" json:"headers"`
}

type Config struct {
	Log            app_log.LogCfg            `envPrefix:"LOG_" json:"log"`
	Application    Application               `envPrefix:"APP_" json:"application"`
	Database       Database                  `envPrefix:"DB_" json:"db"`
	CORS           CORS                      `envPrefix:"CORS_" json:"cors"`
	GuideWebClient http_client.HttpClientCfg `envPrefix:"GUIDE_WEB_CLIENT_" json:"guideWebClient"`
	Bussiness      Bussiness                 `envPrefix:"BUSSINESS_" json:"bussiness"`
}

var (
	instance *Config
	once     sync.Once
)

// Get returns a singleton config loaded from environment variables
func Get() Config {
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
	return *instance
}
