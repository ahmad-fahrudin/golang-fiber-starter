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

# Fresh migration commands (similar to Laravel's migrate:fresh --seed)
migrate-fresh:
	@echo "Running fresh migration (drop all tables and re-migrate)..."
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path src/database/migrations down -all
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path src/database/migrations up
	@echo "Fresh migration completed!"

migrate-fresh-seed:
	@echo "Running fresh migration with seeding (similar to Laravel's migrate:fresh --seed)..."
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path src/database/migrations down -all
	@migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path src/database/migrations up
	@echo "Migration completed, now running seeders..."
	@go run src/main.go --seed all
	@echo "Fresh migration with seeding completed!"

migrate-docker-fresh:
	@echo "Running fresh migration with Docker..."
	@docker run -v ./src/database/migrations:/migrations --network go-fiber-boilerplate_go-network migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable down -all
	@docker run -v ./src/database/migrations:/migrations --network go-fiber-boilerplate_go-network migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up
	@echo "Docker fresh migration completed!"

migrate-docker-fresh-seed:
	@echo "Running fresh migration with seeding using Docker..."
	@docker run -v ./src/database/migrations:/migrations --network go-fiber-boilerplate_go-network migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable down -all
	@docker run -v ./src/database/migrations:/migrations --network go-fiber-boilerplate_go-network migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up
	@echo "Migration completed, now running seeders..."
	@go run src/main.go --seed all
	@echo "Docker fresh migration with seeding completed!"

# Seeder commands
seed-all:
	@go run src/main.go --seed all
seed-list:
	@go run src/main.go --seed list
seed-%:
	@go run src/main.go --seed run $(shell echo $* | sed 's/_/ /g')
seed-refresh-%:
	@go run src/main.go --seed refresh $(word 1,$(subst _, ,$*)) $(word 2,$(subst _, ,$*))
seed-truncate-%:
	@go run src/main.go --seed truncate $*