package main

import (
	"7laba/tspserver/internal"
	"flag"
	"log"
)

func main() {
	addr := flag.String("addr", ":8080", "Сервер для прослушивания")
	flag.Parse()
	if err := internal.NewServer(*addr).Run(); err != nil {
		log.Fatal(err)
	}

}
