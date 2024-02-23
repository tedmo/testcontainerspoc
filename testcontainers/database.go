package testcontainers

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

type TestDatabaseConfig struct {
	// The database username
	Username string
	// The database password
	Password string
	// The database name
	Database string
	// The absolute path of the liquibase/migrations directory.
	// This directory should include the liquibase root changelog and the liquibase Dockerfile.
	MigrationsPath string
}

type TestDatabase struct {
	Host               string
	Port               int
	Username           string
	Password           string
	Database           string
	network            *testcontainers.DockerNetwork
	postgresContainer  testcontainers.Container
	liquibaseContainer testcontainers.Container
}

// WithNetwork is a helper function to allow setting the Network option in the testcontainers postgres module
func WithNetwork(network string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Networks = []string{
			network,
		}
	}
}

// WithName is a helper function to allow setting the Name option in the testcontainers postgres module
func WithName(name string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Name = name
	}
}

// NewTestDatabase launches a postgres instance, and executes schema migrations using liquibase.
// The caller should call Close when finished to shut it down.
//
// TODO: Should this return errors or fail the tests?  Need to weigh pros/cons
func NewTestDatabase(t *testing.T, config *TestDatabaseConfig) *TestDatabase {
	t.Helper()

	ctx := context.Background()

	// Creating a network for both containers (postgres and liquibase) to run in so we can ensure the containers can
	// communicate with each other.  During initial implementation, it was attempted to run postgres on the default
	// bridge network, and set the liquibase container network to "host", but it never quite worked.  The only thing
	// that worked was creating a shared network and setting liquibase to connect to the container IP and container
	// port (5432) like we're doing here.
	// TODO: Look into this further to see if can avoid this (and how that impacts running on CI servers)

	// Network
	dbNetwork, err := network.New(ctx,
		network.WithCheckDuplicate(),
		network.WithAttachable(),
		network.WithDriver("bridge"),
	)
	require.NoError(t, err, "failed to create container network")

	// Postgres
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase(config.Database),
		postgres.WithUsername(config.Username),
		postgres.WithPassword(config.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
		WithNetwork(dbNetwork.Name),
		WithName("db"),
	)
	t.Cleanup(func() {
		postgresContainer.Terminate(ctx)
	})
	require.NoError(t, err, "failed to start postgres container")

	dbHost, err := postgresContainer.Host(ctx)
	require.NoError(t, err, "error getting postgres container host")

	dbPort, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err, "error getting postgres container port")

	postgresIP, err := postgresContainer.ContainerIP(ctx)
	require.NoError(t, err, "error getting postgres container ip")

	// Liquibase Migrations
	postgresURL := fmt.Sprintf("jdbc:postgresql://%s:%d/%s?sslmode=disable", postgresIP, 5432, config.Database)

	// Building the liquibase container from a Dockerfile that extends the official liquibase image and copies the
	// necessary files (changelog + DDL scripts) into the image. This is done because we cannot bind to a directory
	// with testcontainers.
	// TODO: Find a way to use testcontainers `Files` property to copy the necessary files so we can avoid building
	//   a customer image.
	liquibaseContainerReq := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    config.MigrationsPath,
			Dockerfile: "Dockerfile",
			KeepImage:  false,
		},
		Cmd: []string{
			"update",
			"--defaults-file", "/liquibase/changelog/liquibase.docker.properties",
			"--url", postgresURL,
			"--username", config.Username,
			"--password", config.Password,
			"--changelog-file", "changelog.xml",
		},
		Networks:   []string{dbNetwork.Name},
		Name:       "migrations",
		WaitingFor: wait.ForExit().WithExitTimeout(10 * time.Second),
	}

	liquibaseContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: liquibaseContainerReq,
		Started:          true,
	})
	require.NoError(t, err, "failed to start liquibase container")

	return &TestDatabase{
		Host:               dbHost,
		Port:               dbPort.Int(),
		Username:           config.Username,
		Password:           config.Password,
		Database:           config.Database,
		network:            dbNetwork,
		postgresContainer:  postgresContainer,
		liquibaseContainer: liquibaseContainer,
	}
}

func (db *TestDatabase) Close(t *testing.T, ctx context.Context) {
	t.Helper()

	if db.liquibaseContainer != nil && db.liquibaseContainer.IsRunning() {
		if err := db.liquibaseContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to termination liquibase container: %v\n", err)
		}
	}

	if db.postgresContainer != nil && db.postgresContainer.IsRunning() {
		if err := db.postgresContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to termination postgres container: %v\n", err)
		}
	}

	if db.network != nil {
		if err := db.network.Remove(ctx); err != nil {
			fmt.Printf("failed to remove testcontainers network: %v\n", err)
		}
	}
}
