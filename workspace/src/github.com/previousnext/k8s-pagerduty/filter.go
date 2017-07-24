package main

import (
	"fmt"
)

const (
	labelNamespace = "io.kubernetes.pod.namespace"
	labelPod       = "io.kubernetes.pod.name"
	labelName      = "io.kubernetes.container.name"
)

// Helper function to allow us to filter down which containers we should be pushing to our
// CloudWatch Logs backend.
func filter(labels map[string]string) (string, string, string, error) {
	if _, ok := labels[labelNamespace]; !ok {
		return "", "", "", fmt.Errorf("cannot find namespace label")
	}

	if _, ok := labels[labelPod]; !ok {
		return "", "", "", fmt.Errorf("cannot find pod label")
	}

	if _, ok := labels[labelName]; !ok {
		return "", "", "", fmt.Errorf("cannot find name label")
	}

	// Don't log container which are just in place to manage the Kubernetes POD.
	if labels[labelName] == "POD" {
		return "", "", "", fmt.Errorf("this is a POD container")
	}

	return labels[labelNamespace], labels[labelPod], labels[labelName], nil
}
