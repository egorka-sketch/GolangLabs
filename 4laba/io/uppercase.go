package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Uppercase() {
	fmt.Print("Задание 4\nВведите строку: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	line = strings.TrimSpace(line)
	fmt.Println("Строка в верхнем регистре: ", strings.ToUpper(line))
}
