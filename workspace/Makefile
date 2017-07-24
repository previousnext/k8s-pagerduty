#!/usr/bin/make -f

build:
	./hack/build.sh linux server k8s-pagerduty github.com/previousnext/k8s-pagerduty

lint:
	golint -set_exit_status src/github.com/previousnext/k8s-pagerduty/...

test:
	./hack/test.sh github.com/previousnext/k8s-pagerduty/...
