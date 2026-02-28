package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/di"
)

type Server struct {
	router http.Handler
	srv    *http.Server
}

func NewServer(router http.Handler) *Server {
	return &Server{
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	address := fmt.Sprintf(":%d", di.ProvideConfiguration().Port())
	log.Println("Starting server on address:", address)

	s.srv = &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.srv.Shutdown(ctx)
}
