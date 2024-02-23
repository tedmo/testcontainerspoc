include .env

### App
.PHONY: go/run
go/run: db/schema/update
	DB_USER=${DB_USER} \
	DB_PASSWORD=${DB_PASSWORD} \
	DB_HOST=${DB_HOST} \
	DB_PORT=${DB_PORT} \
	DB_NAME=${DB_NAME} \
	go run cmd/app/main.go

### Go
.PHONY: go/update
go/update:
	go get -u -t ./...
	go mod tidy

### SQLC
.PHONY: sqlc/gen
sqlc/gen:
	sqlc generate

### Database
.PHONY: db/up
db/up:
	docker compose up --detach

.PHONY: db/down
db/down:
	docker compose down

### Schema
.PHONY: db/schema/update
db/schema/update: db/up
	docker run --network=host --rm -v ${CURDIR}/migrations:/liquibase/changelog liquibase/liquibase:4.26 update-testing-rollback \
	  --url jdbc:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME} \
	  --username ${DB_USER} \
	  --password ${DB_PASSWORD} \
	  --changelog-file=changelog.xml \
	  --defaults-file=/liquibase/changelog/liquibase.docker.properties

