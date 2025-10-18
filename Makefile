.PHONY: build clean run

build:
	@mkdir -p bin
	@go build -o bin/clipd ./cmd/clipd
	@go build -o bin/clipctl ./cmd/clipctl

clean:
	@rm -rf bin

run: build
	@./bin/clipd