package main

import (
	"10laba/tlsclient/internal"
	"flag"
	"log"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "Адрес сервера")
	flag.Parse()
	if err := internal.SendMsg(*addr); err != nil {
		log.Fatal(err)
	}

}
