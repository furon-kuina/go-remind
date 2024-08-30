package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"golang.org/x/sync/errgroup"
)

// config holds every information needed to run the application
type config struct {
	port      int    `env:"PORT"`
	redisHost string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	redisPort int    `env:"REDIS_PORT" envDefault:"36379"`
}

// main checks whether given config is valid, and is responsible for
// exiting the program with the appropriate exit status
func main() {
	cfg := &config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Printf("invalid config: %v", err)
		os.Exit(1)
	}
	if err := run(cfg); err != nil {
		log.Printf("server finished with error: %v", err)
		os.Exit(1)
	}
}

// run runs all the goroutines, and return value of run
// defines the program's status code
func run(cfg *config) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.port))
	if err != nil {
		return err
	}
	log.Printf("starting server at http://%s", l.Addr().String())
	mux := NewMux()
	s := NewServer(l, mux)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.s.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close server: %v", err)
			return err
		}
		return nil
	})
	<-ctx.Done()
	if err := s.s.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown failed: %v", err)
	}
	return eg.Wait()
}

// Server contains every every information needed to run the server
type Server struct {
	s *http.Server
	l net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		s: &http.Server{Handler: mux},
		l: l,
	}
}
