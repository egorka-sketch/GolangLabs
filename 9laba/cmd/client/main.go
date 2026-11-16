package main

import (
	"flag"
	"os"

	"9laba/internal/client/api"
	"9laba/internal/client/load"
	"9laba/internal/client/ui"
)

func main() {
	addr := flag.String("addr", "http://localhost:8080", "адрес сервера")
	workers := flag.Int("load", 0, "режим нагрузки: количество одновременных клиентов")
	iterations := flag.Int("iter", 10, "итераций на каждого клиента в режиме нагрузки")
	login := flag.String("login", "", "логин для режима нагрузки")
	password := flag.String("password", "", "пароль для режима нагрузки")
	flag.Parse()

	if *workers > 0 {
		load.Run(*addr, *login, *password, *workers, *iterations)
		return
	}
	ui.NewMenu(api.New(*addr), os.Stdin).Run()
}
