package alfred

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLimitSlice(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	assert.Equal(t, 2, len(LimitSlice(arr, 2)))
	assert.Equal(t, 0, len(LimitSlice(arr, 0)))
}
