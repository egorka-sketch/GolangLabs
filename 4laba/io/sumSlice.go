package io

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func SumSlice() {
	sum := 0
	fmt.Print("Задание 5\nВведите числа через пробел: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	line = strings.TrimSpace(line)
	elements := strings.Fields(line)
	for _, element := range elements {
		strconv.Atoi(element)
		num, err := strconv.Atoi(element)
		if err != nil {
			fmt.Println("Не является числом", element)
			continue
		}
		sum += num
	}
	fmt.Println("Сумма чисел:", sum)
}
