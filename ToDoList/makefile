DB_DSN := "postgres://postgres:12345@localhost:5438/postgres?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)
gen:
	oapi-codegen -config openapi/.openapi -include-tags tasks -package tasks openapi/openapi.yaml > ./internal/web/tasks/api.gen.go
migrate-f: 
	$(MIGRATE) force $V

migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}



migrate:
	$(MIGRATE) up


migrate-down:
	$(MIGRATE) down
	

run:
	go run cmd/app/main.go 		