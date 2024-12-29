build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run:
	@./bin/ecom

migration:
	@migration create -ext sql -dir cmd/migrate/migration $(filter-out $@,$(MAKECMDGOALS))

migration-up:
	@go run cmd/migrate/main.go up

migration-down:
	@go run cmd/migrate/main.go down
