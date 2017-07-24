package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	var store Store

	// Test values that are "too high".
	err := store.Add(101)
	if assert.NotNil(t, err) {
		assert.Equal(t, "value is greater than 100", err.Error())
	}

	// Test values that are "too low".
	err = store.Add(-1)
	if assert.NotNil(t, err) {
		assert.Equal(t, "value is less than 0", err.Error())
	}

	err = store.Add(100)
	assert.Nil(t, err)
	_, err = store.Avg()
	if assert.NotNil(t, err) {
		assert.Equal(t, "second data point does not contain data (or zero value)", err.Error())
	}

	err = store.Add(100)
	assert.Nil(t, err)
	_, err = store.Avg()
	if assert.NotNil(t, err) {
		assert.Equal(t, "third data point does not contain data (or zero value)", err.Error())
	}

	err = store.Add(100)
	assert.Nil(t, err)
	avg, err := store.Avg()
	assert.Nil(t, err)
	assert.Equal(t, 100, avg)

	err = store.Add(20)
	assert.Nil(t, err)
	avg, err = store.Avg()
	assert.Nil(t, err)
	assert.Equal(t, 73, avg)

	err = store.Add(30)
	assert.Nil(t, err)
	avg, err = store.Avg()
	assert.Nil(t, err)
	assert.Equal(t, 50, avg)

	err = store.Add(50)
	assert.Nil(t, err)
	avg, err = store.Avg()
	assert.Nil(t, err)
	assert.Equal(t, 33, avg)
}
