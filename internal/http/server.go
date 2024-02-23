package http

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/flowchartsman/swaggerui"
	"github.com/tedmo/testcontainerspoc/internal/app"
	"net/http"
	"strconv"
)

//go:embed docs/openapi.yaml
var OpenAPISpec []byte

type UserRepo interface {
	CreateUser(ctx context.Context, user *app.CreateUserReq) (*app.User, error)
	FindUserByID(ctx context.Context, id int64) (*app.User, error)
	FindUsers(ctx context.Context) ([]app.User, error)
}

type Server struct {
	UserRepo UserRepo
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", s.HandleCreateUser())
	mux.HandleFunc("GET /users", s.HandleGetUsers())
	mux.HandleFunc("GET /users/{id}", s.HandleGetUser())
	mux.Handle("/docs/", http.StripPrefix("/docs", swaggerui.Handler(OpenAPISpec)))

	return mux
}

func (s *Server) pathValueInt64(r *http.Request, name string) (int64, error) {
	val := r.PathValue(name)
	int64Val, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64Val, nil
}

func (s *Server) ok(w http.ResponseWriter, v interface{}) {
	s.json(w, http.StatusOK, v)
}

func (s *Server) created(w http.ResponseWriter, v interface{}) {
	s.json(w, http.StatusCreated, v)
}

func (s *Server) json(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) badRequest(w http.ResponseWriter) {
	s.error(w, http.StatusBadRequest, "bad request")
}

func (s *Server) notFound(w http.ResponseWriter) {
	s.error(w, http.StatusNotFound, "not found")
}

func (s *Server) internalError(w http.ResponseWriter) {
	s.error(w, http.StatusInternalServerError, "unexpected error")
}

func (s *Server) error(w http.ResponseWriter, status int, msg string) {
	s.json(w, status, ErrorResponse{Error: msg})
}
