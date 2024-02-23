package main_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tedmo/testcontainerspoc/internal/app"
	"github.com/tedmo/testcontainerspoc/internal/http"
	"github.com/tedmo/testcontainerspoc/internal/postgres"
	"github.com/tedmo/testcontainerspoc/testcontainers"
	nethttp "net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {

	testServer := NewTestServer(t)

	testClient := testServer.Client()
	baseURL := testServer.URL

	t.Cleanup(func() {
		testServer.Close()
	})

	t.Run("create user", func(t *testing.T) {
		reqBody := `{"name": "test"}`
		req, err := nethttp.NewRequest(nethttp.MethodPost, baseURL+"/users", strings.NewReader(reqBody))
		require.NoError(t, err)

		resp, err := testClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, nethttp.StatusCreated, resp.StatusCode)

		var user app.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test", user.Name)
	})

	t.Run("find user by id", func(t *testing.T) {
		req, err := nethttp.NewRequest(nethttp.MethodGet, baseURL+"/users/1", nethttp.NoBody)
		require.NoError(t, err)

		resp, err := testClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, nethttp.StatusOK, resp.StatusCode)

		var user app.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test", user.Name)
	})

	t.Run("find users", func(t *testing.T) {
		req, err := nethttp.NewRequest(nethttp.MethodGet, baseURL+"/users", nethttp.NoBody)
		require.NoError(t, err)

		resp, err := testClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, nethttp.StatusOK, resp.StatusCode)

		var users []app.User
		err = json.NewDecoder(resp.Body).Decode(&users)
		require.NoError(t, err)

		require.Len(t, users, 1)
		user := users[0]
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test", user.Name)
	})

}

func NewTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	ctx := context.Background()

	testDB := testcontainers.NewTestDatabase(t, &testcontainers.TestDatabaseConfig{
		Username:       "user",
		Password:       "password",
		Database:       "users",
		MigrationsPath: filepath.Join(app.RootPath(), "migrations"),
	})
	t.Cleanup(func() {
		testDB.Close(t, ctx)
	})

	db, err := postgres.NewDB(&postgres.DBConfig{
		Host:     testDB.Host,
		Port:     testDB.Port,
		User:     testDB.Username,
		Password: testDB.Password,
		Database: testDB.Database,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})

	server := &http.Server{UserRepo: postgres.NewUserRepo(db)}

	return httptest.NewServer(server.Routes())
}
