package postgres_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tedmo/testcontainerspoc/internal/app"
	pg "github.com/tedmo/testcontainerspoc/internal/postgres"
	"github.com/tedmo/testcontainerspoc/testcontainers"
	"path/filepath"
	"testing"
)

func TestUserRepo(t *testing.T) {

	migrationsPath, err := filepath.Abs(filepath.Join(app.RootPath(), "migrations"))
	require.NoError(t, err, "error getting migrations file path")

	ctx := context.Background()
	testDatabase := testcontainers.NewTestDatabase(t, &testcontainers.TestDatabaseConfig{
		Username:       "user",
		Password:       "pass",
		Database:       "users",
		MigrationsPath: migrationsPath,
	})
	t.Cleanup(func() {
		testDatabase.Close(t, ctx)
	})

	// DB Connection
	db, err := pg.NewDB(&pg.DBConfig{
		Host:     testDatabase.Host,
		Port:     testDatabase.Port,
		User:     testDatabase.Username,
		Password: testDatabase.Password,
		Database: testDatabase.Database,
	})
	t.Cleanup(func() {
		db.Close()
	})

	err = db.Ping()
	require.NoError(t, err)

	// Repo
	repo := pg.NewUserRepo(db)

	var user *app.User
	t.Run("create user", func(t *testing.T) {
		// Test
		user, err = repo.CreateUser(ctx, &app.CreateUserReq{
			Name: "test",
		})
		require.NoError(t, err)
	})

	t.Run("find user by id", func(t *testing.T) {
		foundUser, err := repo.FindUserByID(ctx, user.ID)
		require.NoError(t, err)

		assert.Equal(t, user, foundUser)
	})

	t.Run("find users", func(t *testing.T) {
		foundUsers, err := repo.FindUsers(ctx)
		require.NoError(t, err)

		require.Len(t, foundUsers, 1)
		assert.Equal(t, user, &foundUsers[0])
	})
}
