migrate:
	go run ./cmd/migrator --storage-path=./storage/accounts.db --migrations-path=./migrations

run:
	go run ./cmd/accounts --config=./config/config.yaml
