package db_poll

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"via/internal/config"
	"via/internal/secret"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

var (
	instance *sql.DB
	once     sync.Once
)

func New(cfg config.Database) (*sql.DB, error) {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, secret.ReadSecret(cfg.PasswordFile), cfg.Host, cfg.Port, cfg.Base)
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
