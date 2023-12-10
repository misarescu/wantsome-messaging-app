package server

import (
	"chat-app/internal/storage"
	"chat-app/pkg/loggers"
	"chat-app/pkg/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	router     *mux.Router
	store      storage.Storage
	broadcast  chan models.RoomMessage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	s := Server{
		listenAddr: listenAddr,
		store:      store,
		broadcast:  make(chan models.RoomMessage),
	}

	s.initRouter()

	return &s
}

var shutdown os.Signal = syscall.SIGUSR1

func (s *Server) RunServer() {
	// http.HandleFunc("/", home)

	server := &http.Server{Addr: s.listenAddr}

	go s.handleBroadcastRoomChat()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		loggers.InfoLogger.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			loggers.ErrorLogger.Printf("%s\n", err)
			stop <- shutdown
		}
	}()

	signal := <-stop
	loggers.InfoLogger.Printf("Shutting down server ...")

	server.Shutdown(nil)

	if signal == shutdown {
		os.Exit(1)
	}
}
