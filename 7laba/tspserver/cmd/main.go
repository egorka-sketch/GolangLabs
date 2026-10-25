package main

import (
	"7laba/tspserver/internal"
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":8080", "Сервер для прослушивания")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := internal.NewServer(*addr).Run(ctx); err != nil {
		log.Fatal(err)
	}
}
