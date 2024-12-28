build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run:
	@./bin/ecom