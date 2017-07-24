#!/bin/bash

export GOPATH=$(pwd)/vendor:$(pwd)

go test -cover $1
