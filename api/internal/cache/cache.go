package cache

import (
	"context"
	"encoding/json"
	"via/internal/ds"
)

type Cache struct {
	ds ds.DS
}

func New(ds ds.DS) *Cache {
	return &Cache{ds: ds}
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
	var strVal string
	switch v := any(value).(type) {
	case string:
		strVal = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		strVal = string(bytes)
	}
	return c.ds.Set(ctx, key, strVal, ttlSeconds)

}

func (c *Cache) Get(ctx context.Context, key string, result any) (bool, error) {
	found, value, err := c.ds.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil // Key not found
	}
	switch resultPtr := result.(type) {
	case *string:
		*resultPtr = value
		return true, nil
	default:
		// Decode JSON into struct
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}
