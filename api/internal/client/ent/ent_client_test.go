package ent_client

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingletonBehavior(t *testing.T) {
	// Reset to ensure clean state
	reset()

	db1 := &sql.DB{} // Mock DB (fake for this test)
	client1 := New(db1)

	// Call New again with different DB, should still return same instance
	db2 := &sql.DB{}
	client2 := New(db2)

	assert.Equal(t, client1, client2, "Singleton should always return same instance")
	assert.NotNil(t, Get(), "Get() should return the initialized client")
	assert.Equal(t, client1, Get(), "Get() should match the singleton instance")
}

func TestReset(t *testing.T) {
	reset()
	assert.Nil(t, Get(), "After reset, Get() should return nil")

	db := &sql.DB{}
	client := New(db)
	assert.NotNil(t, client, "New() should initialize client after reset")
}
