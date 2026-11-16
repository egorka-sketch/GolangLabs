package main

import (
	"flag"
	"os"

	"9laba/internal/client/api"
	"9laba/internal/client/ui"
)

func main() {
	addr := flag.String("addr", "http://localhost:8080", "адрес сервера")
	flag.Parse()

	ui.NewMenu(api.New(*addr), os.Stdin).Run()
}
