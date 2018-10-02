GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./... | grep -Ev '(/vendor/|testutil)')
VERSION = 0.0.1
BUILD_LDFLAGS = "\
					-X main.version=$(VERSION)"

all: ensure generate format lint test vet build

lint:
	gometalinter --vendored-linters --vendor --cyclo-over=15 --deadline=100s \
	             --aggregate --sort=path --disable=megacheck --disable=gas \
	             --exclude='/usr/local/go' \
	             --exclude='fmt.Fprint' .

build: $(GOFILES) generate format
	@mkdir -p bin/local/
	go build -ldflags $(BUILD_LDFLAGS) -o bin/local/cloudops

build/linux:
	@mkdir -p bin/linux_amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(BUILD_LDFLAGS) -o bin/linux_amd64/cloudops

generate: install
	go generate ./...

format:
	goimports -local "github.com/ymgyt/cloudops" -w .

test:
	@go test -cover $(GOPACKAGES)

deps:
	go get -u github.com/golang/dep/cmd/dep
	go install github.com/golang/mock/mockgen
	go get -u golang.org/x/tools/cmd/stringer

ensure:
	dep ensure

clean:
	rm -rf bin/ ||:

install:
	go install ./...

vet:
	go vet ./...

.PHONY: lint build format test deps all install vet build-linux
