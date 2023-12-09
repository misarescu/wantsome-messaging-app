package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type apiHandler func(http.ResponseWriter, *http.Request) error

func makeHandler(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle bad requests
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, err)
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(data)
}

func (s *Server) initRouter() {
	s.router = mux.NewRouter()

	chatRouter := s.router.PathPrefix("/chat").Subrouter()
	chatRouter.HandleFunc("/room/{id}", makeHandler(s.handleRoomChat))

	userRouter := s.router.PathPrefix("/users").Subrouter()
	// userRouter.HandleFunc("/{id}", makeHandler(s.handleGetUserByID)).Methods("GET")
	// userRouter.HandleFunc("/{id}", makeHandler(s.handleRemoveUserById)).Methods("DELETE")
	// userRouter.HandleFunc("/{id}", makeHandler(s.handleUpdateUserById)).Methods("PUT")
	userRouter.HandleFunc("", makeHandler(s.handleGetAllUsers)).Methods("GET")
	// userRouter.HandleFunc("", makeHandler(s.handleUpdateMultipleUsers)).Methods("PUT")

	http.Handle("/", s.router)
}
