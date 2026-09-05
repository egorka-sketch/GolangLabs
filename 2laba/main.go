package main

import (
	"fmt"
	"unicode/utf8"
)

type Rectangle struct {
	width, height float64
}

func (r Rectangle) getArea() float64 {
	return r.width * r.height
}
func setArea() Rectangle {
	var width, height float64
	fmt.Print("Введите ширину и высоту через пробел: ")
	fmt.Scan(&width, &height)
	return Rectangle{
		width:  width,
		height: height,
	}

}

func determinantParity() string {
	var number int
	fmt.Print("Задние 1\nВедите число : ")
	fmt.Scan(&number)
	if number%2 == 0 {
		return "Четное"
	} else {
		return "Нечетное"
	}

}

func determinantNumber() string {
	var number int
	fmt.Print("\nЗадние 2\nВведите число : ")
	fmt.Scan(&number)
	if number > 0 {
		return "Positive"
	} else if number < 0 {
		return "Negative"
	}
	return "Zero"

}

func cycle() []int {
	var numbers []int
	fmt.Print("\nЗадание 3\n")
	for i := 1; i <= 10; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

func determinantLength() {
	var line string
	fmt.Print("\nЗадание 4\nВведите строку : ")
	fmt.Scan(&line)
	fmt.Print("Результат : ")
	length := utf8.RuneCountInString(line)
	fmt.Println(length)
}

func average() float64 {
	var a, b int
	fmt.Print("\nВведите a и b через пробел : ")
	fmt.Scan(&a, &b)
	return float64(a+b) / 2
}

func main() {
	fmt.Print(determinantParity())
	fmt.Print(determinantNumber())
	fmt.Print(cycle())
	determinantLength()
	rectangle := setArea()
	fmt.Print(rectangle.getArea())
	fmt.Print(average())
}
