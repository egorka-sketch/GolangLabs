package io

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func InvertedSlice() {
	fmt.Print("Задание 6\nВведите строку: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}
	elements := strings.Fields(strings.TrimSpace(line))
	numbers := make([]int, 0, len(elements))
	for _, element := range elements {
		num, err := strconv.Atoi(element)
		if err != nil {
			continue
		}
		numbers = append(numbers, num)
	}
	slices.Reverse(numbers)
	fmt.Println("В обратном порядке:", numbers)

}
