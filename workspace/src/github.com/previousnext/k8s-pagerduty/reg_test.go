package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReg(t *testing.T) {
	// Meant to be matched.
	matched, err := reg("prod-*", "prod-backup*", "prod-foo")
	assert.Nil(t, err)
	assert.True(t, matched)

	// Meant to be skipped.
	matched, err = reg("prod-*", "prod-backup*", "dev-foo")
	assert.Nil(t, err)
	assert.False(t, matched)

	// Meant to be skipped.
	matched, err = reg("prod-*", "prod-backup*", "prod-backup-foo")
	assert.Nil(t, err)
	assert.False(t, matched)
}
