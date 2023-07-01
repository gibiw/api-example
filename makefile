.PHONY: migration_status
migration_status:
	goose -dir ./migrations postgres "host=127.0.0.1 port=5432 user=postgres password=Qwerty123 dbname=cars sslmode=disable" status

.PHONY: migration_up
migration_up:
	goose -dir ./migrations postgres "host=127.0.0.1 port=5432 user=postgres password=Qwerty123 dbname=cars sslmode=disable" up

.PHONY: migration_down
migration_down:
	goose -dir ./migrations postgres "host=127.0.0.1 port=5432 user=postgres password=Qwerty123 dbname=cars sslmode=disable" reset		

.PHONY: run
run:
	swag init --parseInternal -g cmd/main.go -q
	go run ./cmd/main.go		

.PHONY: test
test:
	go test -v ./... 	

.PHONY: coverage
coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -func ./coverage.out

.PHONY: lint
test:
	golangci-lint run	

.PHONY: database_up
database:
	docker-compose -f deploy/docker-compose.yml -p api up -d