package server

import (
	"chat-app/pkg/loggers"
	"chat-app/pkg/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type apiHandler func(http.ResponseWriter, *http.Request) error

func makeHandler(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggers.InfoLogger.Printf("Responded with: %+v\n", models.ResponseMessage{Message: "valid object",ErrorStatus: true})
		// handle bad requests
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, models.ResponseMessage{Message: err.Error(),ErrorStatus: true})
		}
	}
}

func makeWsHandler(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle bad requests
		if err := f(w, r); err != nil {
			loggers.ErrorLogger.Printf("error with: %+v\n", models.ResponseMessage{Message: err.Error(),ErrorStatus: true})
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
	chatRouter.HandleFunc("/room/{id}", makeWsHandler(s.handleReadRoomChat))

	userRouter := s.router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/{id}", makeHandler(s.handleGetUserById)).Methods("GET")
	// userRouter.HandleFunc("/{id}", makeHandler(s.handleRemoveUserById)).Methods("DELETE")
	// userRouter.HandleFunc("/{id}", makeHandler(s.handleUpdateUserById)).Methods("PUT")
	userRouter.HandleFunc("", makeHandler(s.handleGetAllUsers)).Methods("GET")
	// userRouter.HandleFunc("", makeHandler(s.handleUpdateMultipleUsers)).Methods("PUT")

	http.Handle("/", s.router)
}
