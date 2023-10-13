migrate_up:
	migrate -source file://migrations -database postgres://admin:password@localhost:5433/gofermart?sslmode=disable up

migrate_down:
	migrate -source file://migrations -database postgres://admin:password@localhost:5433/gofermart?sslmode=disable down