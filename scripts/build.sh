#!/bin/bash

APP_VERSION=$1
if [ -z "${APP_VERSION}" ]; then
    APP_VERSION=development
fi
go build -ldflags "-X main.version=${APP_VERSION}" -o build/check-k8s cmd/*.go