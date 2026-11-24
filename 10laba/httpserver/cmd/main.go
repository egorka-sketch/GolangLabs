package main

import (
	"10laba/httpserver/internal/auth"
	"10laba/httpserver/internal/handler"
	"10laba/httpserver/internal/repository"
	"context"
	"crypto/rand"
	_ "embed"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed login.html
var loginPage []byte

func main() {
	csrf := flag.Bool("csrf", true, "защита от CSRF")
	flag.Parse()

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
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			log.Fatal(err)
		}
		log.Println("сгенерирован случайный секрет")
	}
	jm := auth.NewJWTManager(secret, time.Hour)
	accRepo := repository.NewAccRepo(pool)
	authHr := handler.NewAuthHandler(accRepo, store, jm)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", auth.RequireAuth(jm, store, h.HandleGetAll))
	mux.HandleFunc("GET /users/{id}", auth.RequireAuth(jm, store, h.HandleGetById))
	mux.HandleFunc("POST /users", auth.RequireAuth(jm, store, h.CreateUser))
	mux.HandleFunc("PUT /users/{id}", auth.RequireAuth(jm, store, h.UpdateUser))
	mux.HandleFunc("DELETE /users/{id}", auth.RequireAuth(jm, store, auth.RequireRole("admin", h.DeleteUser)))
	mux.HandleFunc("POST /auth/register", authHr.Registration)
	mux.HandleFunc("POST /auth/login", authHr.Login)
	mux.HandleFunc("POST /auth/logout", auth.RequireAuth(jm, store, authHr.Logout))
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(loginPage)
	})

	var root http.Handler = mux
	if *csrf {
		cop := http.NewCrossOriginProtection()
		cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF заблокирован %s %s (Origin: %q, Sec-Fetch-Site: %q)",
				r.Method, r.URL.Path, r.Header.Get("Origin"), r.Header.Get("Sec-Fetch-Site"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"запрос с чужого сайта отклонён"}`))
		}))
		root = cop.Handler(mux)
		log.Println("Защита от CSRF включена")
	} else {
		log.Println("защита от CSRF выключена")
	}

	server := &http.Server{
		Addr:              ":8080",
		Handler:           root,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
