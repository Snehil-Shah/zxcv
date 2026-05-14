BINARY := zxcv

.PHONY: build test lint vet fmt clean run docs

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/$(BINARY)

run: build
	./bin/$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

lint: vet
	golangci-lint run ./...

clean:
	rm -rf bin/

docs:
	go run ./cmd/docs
