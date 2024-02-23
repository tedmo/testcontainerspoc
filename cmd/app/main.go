package main

import (
	"github.com/tedmo/testcontainerspoc/internal/http"
	"github.com/tedmo/testcontainerspoc/internal/postgres"
	"log"
	nethttp "net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	db, err := postgres.NewDBFromEnv()
	if err != nil {
		return err
	}
	defer db.Close()

	server := &http.Server{UserRepo: postgres.NewUserRepo(db)}

	return nethttp.ListenAndServe(":8080", server.Routes())
}
