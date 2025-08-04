package config

import (
	"log"
	"os"
	"sync"
	"via/internal/auth"
	biz_config "via/internal/biz/config"
	http_client "via/internal/client/http"
	db_pool "via/internal/db/pool"
	"via/internal/ds"
	jwt_key "via/internal/jwt"
	app_log "via/internal/log/app"
	"via/internal/middleware"
	"via/internal/pubsub"
	"via/internal/server"

	"github.com/caarlos0/env/v10"
)

type Application struct {
	Env  string `env:"ENV"  envDefault:"production"    json:"env"`
	Name string `env:"NAME" envDefault:"via"           json:"name"`
	//Port           int    `env:"PORT" envDefault:"8080"          json:"port"`
	RequestTimeout int `env:"REQUEST_TIMEOUT" envDefault:"30" json:"requestTimeout"`
}

type Config struct {
	Log            app_log.LogCfg            `envPrefix:"LOG_" json:"log"`
	Application    Application               `envPrefix:"APP_" json:"application"`
	Database       db_pool.DatabaseCfg       `envPrefix:"DB_" json:"db"`
	CORS           middleware.CORSCfg        `envPrefix:"CORS_" json:"cors"`
	GuideWebClient http_client.HttpClientCfg `envPrefix:"GUIDE_WEB_CLIENT_" json:"guideWebClient"`
	Bussiness      biz_config.BussinessCfg   `envPrefix:"BUSSINESS_" json:"bussiness"`
	OAuth          auth.OAuthConfig          `envPrefix:"OAUTH_" json:"oauth"`
	JWT            jwt_key.JWTConfig         `envPrefix:"JWT_" json:"jwt"`
	DS             ds.DSConfig               `envPrefix:"DS_" json:"ds"`
	RestServer     server.ServerConfig       `envPrefix:"REST_" json:"rest"`
	SSEServer      server.ServerConfig       `envPrefix:"SSE_" json:"sse"`
	PubSub         pubsub.PubSubConfig       `envPrefix:"PUBSUB_" json:"pubsub"`
}

var (
	instance *Config
	once     sync.Once
	mutex    sync.Mutex
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

func reset() {
	mutex.Lock()
	defer mutex.Unlock()
	instance = nil
	once = sync.Once{}
}
