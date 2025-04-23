DB_DSN := "postgres://postgres:12345@localhost:5438/postgres?sslmode=disable"
MIGRATE := migrate -path ./migrations -database "$(DB_DSN)"
MIGRATE := migrate -source file://migrations -database "$(DB_DSN)"


gen:
	oapi-codegen -config openapi/.openapi -include-tags tasks -package tasks openapi/openapi.yaml > ./internal/web/tasks/api.gen.go


gen-users:
	oapi-codegen -config openapi/.openapi -include-tags users -package users openapi/openapi.yaml > ./internal/web/users/api.gen.go

gen-all: gen gen-users
lint:
	golangci-lint run --out-format=colored-line-number


migrate-force:
	$(MIGRATE) force $(V)
 

migrate-new:
	migrate create -ext sql -dir ./migrations $(NAME) 

migrate-up:
	$(MIGRATE) up 
 

migrate-down:
	$(MIGRATE) down 
 


run:
	go run cmd/app/main.go


all: gen gen-users lint migrate-up run 


tag-users:
	git tag -a users/$(shell date +%Y%m%d%H%M%S) -m "Tag for users API"
	git push origin --tags

# Псевдоним для all
default: all