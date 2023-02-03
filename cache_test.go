package alfred

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheData(t *testing.T) {
	data, hit := CacheData("TestCacheData", func() string {
		return "value"
	})
	assert.False(t, hit)
	assert.Equal(t, "value", data)

	data, hit = CacheData("TestCacheData", func() string {
		return "value"
	})
	assert.True(t, hit)
	assert.Equal(t, "value", data)
}

func TestCacheExpiredData(t *testing.T) {
	_, hit := CacheExpiredData("TestCacheExpiredData", 500, func() string {
		return "value"
	})
	assert.False(t, hit)
	time.Sleep(time.Second)

	var loaded bool
	_, hit = CacheExpiredData("TestCacheExpiredData", 500, func() string {
		loaded = true
		return "value"
	})
	assert.False(t, hit)
	assert.True(t, loaded)
}
