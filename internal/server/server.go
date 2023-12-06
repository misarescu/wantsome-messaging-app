package server

import (
	"chat-app/internal/storage"
	"chat-app/pkg/loggers"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	store      storage.Storage
	router     *mux.Router
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	s := Server{
		listenAddr: listenAddr,
		store:      store,
	}

	s.initRouter()

	return &s
}

var shutdown os.Signal = syscall.SIGUSR1

func (s *Server) RunServer() {
	// http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleConnections)

	go handleMsg()

	server := &http.Server{Addr: s.listenAddr}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		loggers.InfoLogger.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			loggers.ErrorLogger.Printf("error starting server %s\n", err)
			stop <- shutdown
		}
	}()

	signal := <-stop
	loggers.InfoLogger.Printf("Shutting down server ...")

	m.Lock()
	for conn := range userConnections {
		conn.Close()
		delete(userConnections, conn)
	}
	m.Unlock()

	server.Shutdown(nil)

	if signal == shutdown {
		os.Exit(1)
	}
}
