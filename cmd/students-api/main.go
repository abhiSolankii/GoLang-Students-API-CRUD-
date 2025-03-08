package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhiSolankii/students-api-go-lang/internal/config"
	"github.com/abhiSolankii/students-api-go-lang/internal/http/handlers/student"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//load database

	//setup routes
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	//setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("address", cfg.HTTPServer.Addr))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	<-done

	//graceful shutdown
	slog.Info("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server ", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")

}
