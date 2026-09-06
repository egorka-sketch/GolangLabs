package main

import (
	"3laba/array"
	"3laba/mathutils"
	"3laba/stringutils"
	"fmt"
)

func main() {

	n := 0
	fmt.Println("Задание 1,2\n Введите n = ")
	fmt.Scan(&n)
	factorial := mathutils.Factorial(n)
	fmt.Println("Факториал = ", factorial)

	var line string
	fmt.Println("Задание 3\n Введите строку: ")
	fmt.Scan(&line)
	turn := stringutils.Reverse(line)
	fmt.Println("Перевернутая строка: ", turn)

	arr := array.GenerateArray()
	fmt.Println("Задание 4\n Массив: ", arr)

	fmt.Print("Задание 5\n")
	array.Slice()

	fmt.Print("Задание 6\n")
	array.SliceWord()
}
