package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	t.Parallel()

	t.Run("should return the value of the environment variable", func(t *testing.T) {
		t.Cleanup(func() {
			os.Clearenv()
		})
		// Set the environment variable
		os.Setenv("TEST_VAR", "test")

		assert.Equal(t, Get("TEST_VAR", "default"), "test")
	})

	t.Run("should return the default value if the environment variable is not set", func(t *testing.T) {
		t.Cleanup(func() {
			os.Clearenv()
		})

		assert.Equal(t, Get("TEST_VAR", "default"), "default")
	})
}
