VERSION = $(shell git describe --tags --abbrev=0)
HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(HASH) -X main.BuildDate=$(DATE)" -o dtgrep-linux-amd64

release: build
	github-release release --user mdom --repo dtgrep --tag $(VERSION)
	github-release upload --user mdom --repo dtgrep --tag $(VERSION) --name dtgrep-linux-amd64 --file dtgrep-linux-amd64

coverage:
	go test -coverprofile=cover.out
	go tool cover -html=cover.out
