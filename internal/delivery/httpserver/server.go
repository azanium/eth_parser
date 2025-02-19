package httpserver

import (
	"context"
	"eth_parser/internal/app/parser"
	"eth_parser/internal/app/repo"
	"eth_parser/internal/delivery/httpserver/middleware"

	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server  *http.Server
	handler *TransactionHandler
	port    string
}

func NewServer(port string) *Server {
	repo := repo.NewMemoryTransactionRepo()
	parser := parser.NewEthereumParser(&http.Client{Timeout: 5 * time.Second}, repo)

	// Initialize handler
	handler := NewTransactionHandler(parser)

	return &Server{
		handler: handler,
		port:    port,
	}
}

func (s *Server) setup() {
	mux := http.NewServeMux()

	mux.HandleFunc("/get-current-block", s.handler.GetCurrentBlock)
	mux.HandleFunc("/subscribe", s.handler.Subscribe)
	mux.HandleFunc("/get-transaction/", s.handler.GetTransaction)

	// Wrap the mux with the recovery middleware
	handler := middleware.Recovery(mux)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Start(errChan chan error) {
	s.setup()
	log.Printf("Server starting on port %s", s.port)

	go func() {
		errChan <- s.server.ListenAndServe()
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
