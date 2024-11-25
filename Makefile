build:
	@go build -o bin/cloudfilex

run: build
	@./bin/cloudfilex

test:
	@go test ./... -v