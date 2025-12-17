migrate-create:
	migrate create -ext sql -dir internal/db/migrations -seq $(NAME)

migrate-up:
	migrate -database $(DATABASE_URL) -path internal/db/migrations up

migrate-down:
	migrate -database $(DATABASE_URL) -path internal/db/migrations down