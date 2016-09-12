VERSION = $(shell git describe --tags --abbrev=0)
HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

PACKAGES = $(shell go list ./...)

build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(HASH) -X main.BuildDate=$(DATE)" -o dtgrep-linux-amd64

release: build
	github-release release --user mdom --repo dtgrep --tag $(VERSION)
	github-release upload --user mdom --repo dtgrep --tag $(VERSION) --name dtgrep-linux-amd64 --file dtgrep-linux-amd64

cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

check:
	go vet ./...
	@echo gofmt -s -l .
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"

.PHONY: check cover release build
