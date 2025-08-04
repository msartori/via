package db_pool

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"via/internal/secret"
	mock_secret "via/internal/secret/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew_Success(t *testing.T) {
	reset()

	secretMock := new(mock_secret.MockSecret)
	secret.Set(secretMock)
	secretMock.On("Read", mock.Anything).Return("test-password")

	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	// Inject mock DB via hook by overriding sql.Open temporarily
	originalSQLOpen := openDB
	openDB = func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}
	defer func() { openDB = originalSQLOpen }()

	mock.ExpectPing().WillReturnError(nil)

	cfg := DatabaseCfg{
		PasswordFile: "secret",
		User:         "user",
		Base:         "base",
		Port:         "5432",
		Host:         "localhost",
	}

	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, db, Get())

	// Singleton check: calling New again must reuse same db
	db2, err := New(cfg)
	assert.NoError(t, err)
	assert.Equal(t, db, db2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNew_ErrorOpeningDB(t *testing.T) {
	reset()

	secretMock := new(mock_secret.MockSecret)
	secret.Set(secretMock)
	secretMock.On("Read", mock.Anything).Return("test-password")

	originalSQLOpen := openDB
	openDB = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, fmt.Errorf("db open failed")
	}
	defer func() { openDB = originalSQLOpen }()

	cfg := DatabaseCfg{}
	db, err := New(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "opening DB")
}

func TestNew_ErrorPingingDB(t *testing.T) {
	reset()

	secretMock := new(mock_secret.MockSecret)
	secret.Set(secretMock)
	secretMock.On("Read", mock.Anything).Return("test-password")

	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	originalSQLOpen := openDB
	openDB = func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}
	defer func() { openDB = originalSQLOpen }()

	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	cfg := DatabaseCfg{}
	db, err := New(cfg)
	assert.Error(t, err)
	assert.NotNil(t, db)
	assert.Contains(t, err.Error(), "pinging DB")
}

func TestReset(t *testing.T) {
	reset()

	secretMock := new(mock_secret.MockSecret)
	secret.Set(secretMock)
	secretMock.On("Read", mock.Anything).Return("test-password")

	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	originalSQLOpen := openDB
	openDB = func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}
	defer func() { openDB = originalSQLOpen }()

	mock.ExpectPing().WillReturnError(nil)

	cfg := DatabaseCfg{}
	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	reset()
	assert.Nil(t, Get())
}
