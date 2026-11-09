package main

import (
	"8laba/internal/handler"
	"8laba/internal/repository"
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://lab8:lab8@localhost:5432/users?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := repository.NewPool(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	repo := repository.NewRepo(pool)
	h := handler.NewHandler(repo)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", h.HandleGetAll)
	mux.HandleFunc("GET /users/{id}", h.HandleGetById)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
