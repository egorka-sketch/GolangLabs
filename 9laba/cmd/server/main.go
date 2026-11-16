package main

import (
	"9laba/internal/server/auth"
	"9laba/internal/server/handler"
	"9laba/internal/server/repository"
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
	store := auth.NewSessionStore()
	accRepo := repository.NewAccRepo(pool)
	authHr := handler.NewAuthHandler(accRepo, store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", auth.RequireAuth(store, h.HandleGetAll))
	mux.HandleFunc("GET /users/{id}", auth.RequireAuth(store, h.HandleGetById))
	mux.HandleFunc("POST /users", auth.RequireAuth(store, h.CreateUser))
	mux.HandleFunc("PUT /users/{id}", auth.RequireAuth(store, h.UpdateUser))
	mux.HandleFunc("DELETE /users/{id}", auth.RequireAuth(store, h.DeleteUser))
	mux.HandleFunc("POST /auth/register", authHr.Registration)
	mux.HandleFunc("POST /auth/login", authHr.Login)
	mux.HandleFunc("POST /auth/logout", auth.RequireAuth(store, authHr.Logout))
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
