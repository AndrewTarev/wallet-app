create-migrations:
	migrate create -ext sql -dir ./internal/infrastructure/database/migrations -seq init

migrateup:
	migrate -path ./internal/infrastructure/database/migrations/ -database 'postgres://postgres:postgres@localhost:5432/wallet-app?sslmode=disable' up

migratedown:
	migrate -path ./internal/infrastructure/database/migrations/ -database 'postgres://postgres:postgres@localhost:5432/wallet-app?sslmode=disable' down

test-mock:
	mockgen -source=internal/app/services/service.go -destination=internal/app/services/mocks/mock_wallet_service.go -package=mocks

gen-docs:
	swag init -g ./cmd/main.go -o ./docs