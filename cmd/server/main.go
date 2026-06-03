package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sssnukke/url-shortener/internal/handler"
	"github.com/sssnukke/url-shortener/internal/repository"
	"github.com/sssnukke/url-shortener/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	ctx := context.Background()

	log.Println("running migrations...")
	if err := repository.RunMigrations(os.Getenv("POSTGRES_DSN")); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations done")

	pool, err := repository.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()
	log.Println("connected to postgres")

	redis, err := repository.NewRedisPool(ctx)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redis.Close()
	log.Println("connected to redis")

	urlRepo := repository.NewPostgresRepo(pool)
	cacheRepo := repository.NewRedisRepo(redis)

	svc := service.NewShortenerService(urlRepo, cacheRepo)

	router := handler.NewRouter(svc)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server listening on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped")
}
