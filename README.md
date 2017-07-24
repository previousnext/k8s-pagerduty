Kubernetes: PagerDuty
=====================

[![CircleCI](https://circleci.com/gh/previousnext/k8s-pagerduty.svg?style=svg)](https://circleci.com/gh/previousnext/k8s-pagerduty)

DaemonSet for tracking container CPU and Memory usage.

A PagerDuty event will be created if a container goes above a CPU or Memory threshold.

## Usage

```bash
kubectl create -f kubernetes/daemonset.yaml
```

## Development

```bash
# Run the test suite
cd workspace && make lint
cd workspace && make test

# Build the project
cd workspace && make build

# Release the project
make release VERSION=1.0.0
```
