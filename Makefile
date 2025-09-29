include .env
export $(shell sed 's/=.*//' .env)

start:
	@go run src/main.go
lint:
	@golangci-lint run
tests:
	@go test -v ./test/...
tests-%:
	@go test -v ./test/... -run=$(shell echo $* | sed 's/_/./g')
testsum:
	@cd test && gotestsum --format testname
swagger:
	@cd src && swag init
migration-%:
	@migrate create -ext sql -dir src/database/migrations create-table-$(subst :,_,$*)
migrate-up:
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path src/database/migrations up
migrate-down:
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path src/database/migrations down
migrate-docker-up:
	@docker run -v ./src/database/migrations:/migrations --network go-fiber-boilerplate_go-network migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up
migrate-docker-down:
	@docker run -v ./src/database/migrations:/migrations --network go-fiber-boilerplate_go-network migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable down -all
docker:
	@chmod -R 755 ./src/database/init
	@docker-compose up --build
docker-test:
	@docker-compose up -d && make tests
docker-down:
	@docker-compose down --rmi all --volumes --remove-orphans
docker-cache:
	@docker builder prune -f

# Seeder commands
seed-all:
	@go run cmd/seeder/main.go all
seed-list:
	@go run cmd/seeder/main.go list
seed-%:
	@go run cmd/seeder/main.go run $(shell echo $* | sed 's/_/ /g')
seed-refresh-%:
	@go run cmd/seeder/main.go refresh $(word 1,$(subst _, ,$*)) $(word 2,$(subst _, ,$*))
seed-truncate-%:
	@go run cmd/seeder/main.go truncate $*