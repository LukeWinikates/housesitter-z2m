.PHONY: test run

build/housesitter-z2m-linux-amd64: $(shell find . -iname "*.go")
	mkdir -p build/
	GOOS=linux GOARCH=amd64 \
		go build  -o $@ cmd/main.go

test:
	go fmt ./...
	go test -v ./...
	go vet ./...
	golangci-lint run

run:
	go run cmd/main.go

