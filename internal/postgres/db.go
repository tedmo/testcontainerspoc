package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewDB(config *DBConfig) (*sql.DB, error) {
	datasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database)

	return sql.Open("postgres", datasource)
}

func NewDBFromEnv() (*sql.DB, error) {
	return NewDB(&DBConfig{
		Host: os.Getenv("DB_HOST"),
		Port: func() int {
			port, err := strconv.Atoi(os.Getenv("DB_PORT"))
			if err != nil {
				return 0
			}
			return port
		}(),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	})
}
