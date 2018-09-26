VERSION = 0.0.1
BUILD_LDFLAGS = "\
					-X main.version=$(VERSION)"
all: ensure generate format lint test build

lint:
	gometalinter --vendored-linters --vendor --cyclo-over=15 \
	             --aggregate --sort=path --disable=megacheck --disable=gas \
	             --exclude='fmt.Fprint'

build: generate format
	@go build -ldflags $(BUILD_LDFLAGS)
	@go install ./...

generate:
	@go generate ./...

format:
	@goimports -local "github.com/ymgyt/cloudops" -w .

test:
	go test -v ./...

deps:
	go get -u github.com/golang/dep/cmd/dep

ensure:
	@dep ensure

clean:
	@rm ./cloudops

.PHONY: lint build format test deps all
