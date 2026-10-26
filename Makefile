.PHONY: test run image

build/housesitter-z2m-linux-amd64: $(shell find . -iname "*.go")
	mkdir -p build/
	GOOS=linux GOARCH=amd64 \
		go build  -o $@ cmd/main.go

image: build/housesitter-z2m-linux-amd64
	docker build -t lukewinikates/housesitter-z2m .

test:
	go fmt ./...
	go test -v ./...
	go vet ./...
	golangci-lint run

run:
	go run cmd/main.go

watch:
	-rm build/server
	MQTT_NO_OP_CLIENT=yes ~/go/bin/air
