package runtime_settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildKey(t *testing.T) {
	key := buildKey("test")
	assert.Equal(t, "settings:test", key)
	key = buildKey(1)
	assert.Equal(t, "settings:1", key)
	key = buildKey(1.1)
	assert.Equal(t, "settings:1.1", key)
	key = buildKey(true)
	assert.Equal(t, "settings:true", key)
	key = buildKey(false)
	assert.Equal(t, "settings:false", key)
	key = buildKey(nil)
	assert.Equal(t, "settings:", key)
}
