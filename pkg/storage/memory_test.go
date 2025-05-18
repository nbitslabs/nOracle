package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	t.Run("Operations", func(t *testing.T) {
		store := NewMemory[string]()
		err := store.Store("key", "value")
		assert.Nil(t, err)

		v, err := store.Get("key")
		assert.Nil(t, err)
		assert.Equal(t, "value", v)

		err = store.Delete("key")
		assert.Nil(t, err)

		_, err = store.Get("key")
		assert.NotNil(t, err)

		err = store.Close()
		assert.Nil(t, err)
	})

	t.Run("Close to ensure delete", func(t *testing.T) {
		store := NewMemory[string]()
		err := store.Store("key", "value")
		assert.Nil(t, err)

		err = store.Close()
		assert.Nil(t, err)
	})

	t.Run("Closed and nil store", func(t *testing.T) {
		store := NewMemory[string]()
		err := store.Close()
		assert.Nil(t, err)

		// Now, every method you call will return an error
		err = store.Store("key", "value")
		assert.NotNil(t, err)

		_, err = store.Get("key")
		assert.NotNil(t, err)

		err = store.Delete("key")
		assert.NotNil(t, err)

		// Except for closing again
		err = store.Close()
		assert.Nil(t, err)
	})
}
