migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/ecom" -path ./cmd/migrate/migrations force 20241230145845

migrate -path ./cmd/migrate/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/ecom" version

migrate -path ./cmd/migrate/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/ecom" force <previous_version>

migrate -path ./cmd/migrate/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/ecom" up 1
