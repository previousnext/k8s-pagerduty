#!/usr/bin/make -f

VERSION=$(shell git describe --tags --always)
IMAGE=previousnext/k8s-pagerduty

release: build push

build:
	docker build -t ${IMAGE}:${VERSION} .

push:
	docker push ${IMAGE}:${VERSION}
