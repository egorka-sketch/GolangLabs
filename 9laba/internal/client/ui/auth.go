package ui

import (
	"fmt"
)

func (m *Menu) register() {
	login, _ := m.prompt("Логин: ")
	password, _ := m.prompt("Пароль: ")
	acc, err := m.client.Register(m.ctx, login, password)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Printf("Аккаунт %s создан, теперь войдите (пункт 7)\n", acc.Login)
}

func (m *Menu) signIn() {
	login, _ := m.prompt("Логин: ")
	password, _ := m.prompt("Пароль: ")
	if err := m.client.Login(m.ctx, login, password); err != nil {
		printErr(err)
		return
	}
	m.login = login
	fmt.Println("Вход выполнен")
}

func (m *Menu) signOut() {
	if !m.client.LoggedIn() {
		fmt.Println("Вы и так не вошли")
		return
	}
	if err := m.client.Logout(m.ctx); err != nil {
		printErr(err)
	}
	m.login = ""
	fmt.Println("Вы вышли из аккаунта")
}
