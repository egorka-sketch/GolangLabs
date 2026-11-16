package ui

import (
	"errors"
	"fmt"
	"net"
	"os"
	"text/tabwriter"
	"time"

	"9laba/internal/client/api"
)

func printUsers(users []api.User) {
	if len(users) == 0 {
		fmt.Println("Ничего не найдено")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tИмя\tДата рождения\tВозраст")
	for _, u := range users {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", u.ID, u.Name, formatDate(u.BirthDay), age(u.BirthDay))
	}
	w.Flush()
	fmt.Printf("Всего на странице: %d\n", len(users))
}

func printUser(u api.User) {
	fmt.Printf("ID:            %d\n", u.ID)
	fmt.Printf("Имя:           %s\n", u.Name)
	fmt.Printf("Дата рождения: %s\n", formatDate(u.BirthDay))
	fmt.Printf("Возраст:       %s\n", age(u.BirthDay))
}

func formatDate(s string) string {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return s
	}
	return t.Format("02.01.2006")
}

func age(s string) string {
	b, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return "?"
	}
	now := time.Now()
	years := now.Year() - b.Year()
	if now.Month() < b.Month() || (now.Month() == b.Month() && now.Day() < b.Day()) {
		years--
	}
	return fmt.Sprintf("%d", years)
}

func printErr(err error) {
	var apiErr *api.APIError
	var netErr net.Error
	switch {
	case errors.As(err, &apiErr):
		fmt.Printf("Ошибка %d: %s\n", apiErr.StatusCode, apiErr.Message)
	case errors.As(err, &netErr) && netErr.Timeout():
		fmt.Println("Сервер не ответил вовремя")
	default:
		fmt.Println("Не удалось связаться с сервером:", err)
	}
}
