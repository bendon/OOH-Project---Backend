build:
	@go build -o bin/production cmd/main.go

test:
	@go test -v ./services/tests/...
	
run: build
	@./bin/production