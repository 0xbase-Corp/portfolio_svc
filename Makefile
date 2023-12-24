migrateup:
	migrate -path pkg/db/migrations -database "postgresql://postgres:postgres@localhost:5555/portfolio?sslmode=disable" -verbose up

migratedown:
	migrate -path pkg/db/migrations -database "postgresql://postgres:postgres@localhost:5555/portfolio?sslmode=disable" -verbose down


.PHONY: migrateup migratedown
