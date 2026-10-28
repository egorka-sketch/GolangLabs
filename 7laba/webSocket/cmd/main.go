package main

import (
	"7laba/webSocket/internal"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	mux := http.NewServeMux()
	center := internal.NewUser()
	mux.HandleFunc("GET /ws", center.ChatHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
		center.CloseAll()
	})

	err := server.ListenAndServe()
	if errors.Is(ctx.Err(), context.Canceled) {
		wg.Wait()
		return
	}
	if err != nil {
		log.Fatal(err)
	}

}
