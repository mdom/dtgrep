VERSION = $(shell git describe --tags --abbrev=0)
HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(HASH) -X main.BuildDate=$(DATE)"
