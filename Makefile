run:
	go run ./cmd/main.go

migrate-up:
	go run ./cmd/migrator/main.go --m=up --storage-path=./storage/image.db --migrations-path=./storage/migrations

migrate-down:
	go run ./cmd/migrator/main.go --m=down --storage-path=./storage/image.db --migrations-path=./storage/migrations