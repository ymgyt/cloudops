GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./... | grep -Ev '(/vendor/|testutil)')
VERSION = 1.0.0
BUILD_LDFLAGS = "\
					-X main.version=$(VERSION)"

all: generate format test vet build

build: $(GOFILES)
	@mkdir -p bin/local/
	go build -ldflags $(BUILD_LDFLAGS) -o bin/local/cloudops

build/linux:
	@mkdir -p bin/linux_amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(BUILD_LDFLAGS) -o bin/linux_amd64/cloudops

generate:
	go generate ./...

format:
	goimports -local "github.com/ymgyt/cloudops" -w .

test:
	@go test -cover $(GOPACKAGES)

deps:
	go install github.com/golang/mock/mockgen
	go get -u golang.org/x/tools/cmd/stringer

clean:
	rm -rf bin/ ||:

vet:
	go vet ./...

.PHONY: lint build format test deps all install vet build-linux
