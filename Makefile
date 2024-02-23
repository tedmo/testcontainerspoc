include .env

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
db/schema/update:
	docker run --network=host --rm -v ${CURDIR}/migrations:/liquibase/changelog liquibase/liquibase:4.26 update-testing-rollback \
	  --url jdbc:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME} \
	  --username ${DB_USER} \
	  --password ${DB_PASSWORD} \
	  --changelog-file=changelog.xml \
	  --defaults-file=/liquibase/changelog/liquibase.docker.properties
