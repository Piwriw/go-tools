LOCAL_VERSION ?=$(shell git rev-parse --abbrev-ref HEAD)_$(shell git config user.email)_$(shell date +"%Y-%m-%d_%H:%M:%S")
CI_COMMIT_SHORT_SHA ?= $(shell git rev-parse --short HEAD)


ifeq ($(shell uname -m),x86_64)
	ARCH := amd64
else
	ARCH := arm64
endif


message:
	@echo ${LOCAL_VERSION}
	@echo ${CI_COMMIT_SHORT_SHA}

build:
     CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build  -gcflags "all=-N -l" -ldflags '-extldflags "-static"' \
     -ldflags "-X ./version.Commit=${CI_COMMIT_SHORT_SHA} -X ./version.Release=${LOCAL_VERSION}" \
     -o main  ./version/version.go;