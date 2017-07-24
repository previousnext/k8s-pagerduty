package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	labels := map[string]string{
		labelNamespace: "test-namespace",
		labelPod:       "test-pod",
		labelName:      "test-name",
	}

	namespace, pod, name, err := filter(labels)
	assert.Nil(t, err)

	assert.Equal(t, "test-namespace", namespace, "Set environment defaults")
	assert.Equal(t, "test-pod", pod, "Set environment defaults")
	assert.Equal(t, "test-name", name, "Set environment defaults")
}

func TestFilterMissingNamespace(t *testing.T) {
	labels := map[string]string{
		labelPod:  "test-pod",
		labelName: "test-name",
	}

	_, _, _, err := filter(labels)
	if assert.NotNil(t, err) {
		assert.Equal(t, "cannot find namespace label", err.Error())
	}
}

func TestFilterMissingPod(t *testing.T) {
	labels := map[string]string{
		labelNamespace: "test-namespace",
		labelName:      "test-name",
	}

	_, _, _, err := filter(labels)
	if assert.NotNil(t, err) {
		assert.Equal(t, "cannot find pod label", err.Error())
	}
}

func TestFilterMissingName(t *testing.T) {
	labels := map[string]string{
		labelNamespace: "test-namespace",
		labelPod:       "test-pod",
	}

	_, _, _, err := filter(labels)
	if assert.NotNil(t, err) {
		assert.Equal(t, "cannot find name label", err.Error())
	}
}

func TestFilterPod(t *testing.T) {
	labels := map[string]string{
		labelNamespace: "test-namespace",
		labelPod:       "test-pod",
		labelName:      "POD",
	}

	_, _, _, err := filter(labels)
	if assert.NotNil(t, err) {
		assert.Equal(t, "this is a POD container", err.Error())
	}
}
