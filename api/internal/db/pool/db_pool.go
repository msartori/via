package db_pool

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"via/internal/secret"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

var (
	instance *sql.DB
	once     sync.Once
)

type DatabaseCfg struct {
	PasswordFile string `env:"PASSWORD_FILE" envDefault:"" json:"passwordFile"`
	Password     string `env:"PASSWORD" envDefault:"" json:"-"`
	User         string `env:"USER" envDefault:"" json:"-"`
	Base         string `env:"BASE" envDefault:"" json:"-"`
	Port         string `env:"PORT" envDefault:"" json:"port"`
	Host         string `env:"HOST" envDefault:"" json:"host"`
}

func New(cfg DatabaseCfg) (*sql.DB, error) {
	var err error
	if cfg.Password == "" {
		cfg.Password = secret.ReadSecret(cfg.PasswordFile)
	}
	once.Do(func() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Base)
		instance, err = sql.Open("pgx", dsn)
		if err != nil {
			err = fmt.Errorf("opening DB: %w", err)
			return
		}
		instance.SetMaxOpenConns(5)
		instance.SetMaxIdleConns(2)
		instance.SetConnMaxLifetime(time.Hour)
		err = instance.Ping()
		if err != nil {
			err = fmt.Errorf("pinging DB: %w", err)
		}
	})
	return instance, err
}

func Get() *sql.DB {
	return instance
}
