package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	store := New()

	err := store.Add("foo", 10, 20)
	assert.Nil(t, err)

	err = store.Add("foo", 10, 20)
	assert.Nil(t, err)

	err = store.Add("foo", 10, 20)
	assert.Nil(t, err)

	cpu, err := store.AvgCPU("foo")
	assert.Nil(t, err)
	assert.Equal(t, 10, cpu)

	mem, err := store.AvgMemory("foo")
	assert.Nil(t, err)
	assert.Equal(t, 20, mem)
}
