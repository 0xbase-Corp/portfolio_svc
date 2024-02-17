migrateup:
	migrate -path pkg/migrations -database "postgresql://postgres:postgres@localhost:5555/portfolio?sslmode=disable" -verbose up

migratedown:
	migrate -path pkg/migrations -database "postgresql://postgres:postgres@localhost:5555/portfolio?sslmode=disable" -verbose down


.PHONY: migrateup migratedown
