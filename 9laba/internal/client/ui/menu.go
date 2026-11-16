package ui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"9laba/internal/client/api"
)

type Menu struct {
	client *api.Client
	in     *bufio.Scanner
	ctx    context.Context
	login  string
}

func NewMenu(c *api.Client, in io.Reader) *Menu {
	return &Menu{client: c, in: bufio.NewScanner(in), ctx: context.Background()}
}

type item struct {
	key, title string
	action     func()
}

func (m *Menu) Run() {
	items := []item{
		{"1", "Список пользователей", m.list},
		{"2", "Найти по ID", m.show},
		{"3", "Добавить", m.create},
		{"4", "Изменить", m.update},
		{"5", "Удалить", m.delete},
		{"6", "Регистрация", m.register},
		{"7", "Вход", m.signIn},
		{"8", "Выход из аккаунта", m.signOut},
	}
	for {
		fmt.Println("\nПользователи")
		if m.client.LoggedIn() {
			fmt.Printf("Вы вошли как %s\n", m.login)
		} else {
			fmt.Println("Вы не вошли")
		}
		for _, it := range items {
			fmt.Printf("%s. %s\n", it.key, it.title)
		}
		fmt.Println("0. Выход")

		choice, ok := m.prompt("Выбор: ")
		if !ok || choice == "0" {
			fmt.Println("До связи")
			return
		}
		found := false
		for _, it := range items {
			if it.key == choice {
				fmt.Println()
				it.action()
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Нет такого пункта")
		}
	}
}

func (m *Menu) list() {
	fmt.Println("Фильтры :")
	var f api.ListFilter
	var ok bool
	if f.Limit, ok = m.optionalInt("Строк на странице : "); !ok {
		return
	}
	if f.Offset, ok = m.optionalInt("Скольких пропустить: "); !ok {
		return
	}
	f.Name, _ = m.prompt("Имя: ")
	if f.MinAge, ok = m.optionalIntPtr("Возраст от: "); !ok {
		return
	}
	if f.MaxAge, ok = m.optionalIntPtr("Возраст до: "); !ok {
		return
	}

	users, err := m.client.ListUsers(m.ctx, f)
	if err != nil {
		printErr(err)
		return
	}
	printUsers(users)
}

func (m *Menu) show() {
	id, ok := m.askID()
	if !ok {
		return
	}
	u, err := m.client.GetUser(m.ctx, id)
	if err != nil {
		printErr(err)
		return
	}
	printUser(u)
}

func (m *Menu) create() {
	name, _ := m.prompt("Имя: ")
	birthDay, _ := m.prompt("Дата рождения (ГГГГ-ММ-ДД): ")
	u, err := m.client.CreateUser(m.ctx, api.UserInput{Name: name, BirthDay: birthDay})
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("Создан пользователь:")
	printUser(u)
}

func (m *Menu) update() {
	id, ok := m.askID()
	if !ok {
		return
	}
	cur, err := m.client.GetUser(m.ctx, id)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("Текущие данные :")
	in := api.UserInput{
		Name:     m.promptDefault("Имя", cur.Name),
		BirthDay: m.promptDefault("Дата рождения", cur.BirthDay),
	}
	u, err := m.client.UpdateUser(m.ctx, id, in)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("Обновлено:")
	printUser(u)
}

func (m *Menu) delete() {
	id, ok := m.askID()
	if !ok {
		return
	}
	answer, _ := m.prompt(fmt.Sprintf("Удалить пользователя %d? (д/н): ", id))
	if a := strings.ToLower(answer); a != "д" {
		fmt.Println("Отменено")
		return
	}
	if err := m.client.DeleteUser(m.ctx, id); err != nil {
		printErr(err)
		return
	}
	fmt.Println("Удалён")
}

func (m *Menu) prompt(label string) (string, bool) {
	fmt.Print(label)
	if !m.in.Scan() {
		return "", false
	}
	return strings.TrimSpace(m.in.Text()), true
}

func (m *Menu) promptDefault(label, current string) string {
	s, _ := m.prompt(fmt.Sprintf("%s [%s]: ", label, current))
	if s == "" {
		return current
	}
	return s
}

func (m *Menu) askID() (int, bool) {
	s, _ := m.prompt("ID: ")
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		fmt.Println("ID должен быть положительным")
		return 0, false
	}
	return id, true
}
func (m *Menu) optionalInt(label string) (int, bool) {
	p, ok := m.optionalIntPtr(label)
	if !ok || p == nil {
		return 0, ok
	}
	return *p, true
}

func (m *Menu) optionalIntPtr(label string) (*int, bool) {
	s, _ := m.prompt(label)
	if s == "" {
		return nil, true
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Нужно число")
		return nil, false
	}
	return &n, true
}
