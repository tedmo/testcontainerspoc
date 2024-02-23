# Testcontainers POC

## Overview
This is a POC experimenting with testcontainers, specifically using testcontainers for bringing up a postgres instance
and running schema migrations with liquibase.

A simple users API with create and read capabilities is used as an example to run some basic database tests and API tests.

## Testcontainers Setup
See `testcontainers/database.go` for the container setup logic.

### Postgres
The postgres DB setup is relatively straight forward using the testingcontainers `postges` module.

### Liquibase Schema Migrations
#### Copying supporting liquibase files to the container
The liquibase setup is a little more complicated.  The official liquibase image documentation suggests using a bind mount
to copy/share the changelog and changeset files with the container, however bind mounts are not supported with testcontainers.
This is why we are building a fresh image that extends the official image and copies the files as part of the image build.

It may be possible to copy the files to the container using the `testcontainers.GenericContainerRequest.Files` field,
but this needs some more investigation.  The issue to overcome with this approach is that we can't copy the files
into a directory that doesn't yet exist in the container, and the liquibase documentation states that the files should
be mounted to the /liquibase/changelog directory, which doesn't exist on the container by default.  The testcontainers 
documentation suggests creating the directory on the container after the container is running, but since we're only using 
the liquibase container to run the liquibase command and exit, that isn't really an option.

#### Communicating between the liquibase container and the postgres container
Another issue was opening communcation between the liquibase and postgres containers.  The eventual solution was to create
a network that both containers run in, and connect liquibase to postgres using the container IP and port, rather than
the host IP and host-exposed port.

With some further experimentation, it may be possible to find a way to use the host IP/port and remove the need for
the network, ultimately reducing complexity.

## Testcontainers Usage Examples

See `internal/postgres/user_test.go` for an example of using the testcontainers postgres instance for testing
the database layer.

See `cmd/main_test.go` for an example of using the testcontainers postgres instance as part of a full integration test
of the API.

## Running the app locally
**Requires Docker installed and running**

This repository includes a `Makefile` with some helpful commands.

These commands are driven from the values set in the `.env` file.

| Command            | Description                                                                                                                   |
|--------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `db/up`            | Bring up a local postgres container (configured in the `docker-compose.yml`) that the app can connect to when running locally |
| `db/down`          | Terminate the postgres container                                                                                              |
| `db/schema/update` | Runs the database schema migrations.  If the database is not already up, the database will be started first.                  |
