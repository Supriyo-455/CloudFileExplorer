build:
	@go build -o bin/distributedStorage

run: build
	@./bin/distributedStorage

test:
	@go test ./... -v