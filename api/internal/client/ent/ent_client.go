package ent_client

import (
	"database/sql"
	"sync"
	"via/internal/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

var (
	instance *ent.Client
	once     sync.Once
	mutex    sync.Mutex
)

func New(dbPool *sql.DB) *ent.Client {
	once.Do(func() {
		driver := entsql.OpenDB(dialect.Postgres, dbPool)
		instance = ent.NewClient(ent.Driver(driver))
	})
	return instance
}

func Get() *ent.Client {
	return instance
}

// Reset resets the singleton (for testing only).
func reset() {
	mutex.Lock()
	defer mutex.Unlock()
	instance = nil
	once = sync.Once{}
}
