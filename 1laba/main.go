package main

import (
	"fmt"
	"time"
)

func Time(t time.Time) string {

	return t.Format("2006-01-02 15:04:05")

}

func Variables() (float64, string, bool) {
	a := 64.35
	b := "eqfqdfw"
	isTrue := false
	return a, b, isTrue
}

func CalculateInt(a, b int) {
	fmt.Printf("a + b =  %d\n", a+b)
	fmt.Printf("a - b =  %d\n", a-b)
	fmt.Printf("a * b =  %d\n", a*b)
	if a%2 == 0 && b%2 == 0 {
		fmt.Printf("a / b =  %d\n", a/b)
	} else {
		fmt.Print("Не получится поделить нечетные числа")
	}
}
func CalculateFloat(a, b float64) {

	fmt.Printf("a + b =  %.2f\n", a+b)
	fmt.Printf("a - b =  %.2f\n", a-b)

}

func CalculateCount(a, b, c float64) {
	fmt.Printf("Среднее значение:  %.3f\n", (a+b+c)/3)
}

func main() {
	fmt.Print("Задание 1\n")
	fmt.Println(Time(time.Now()))
	fmt.Print("\nЗадание 2 и 3\n")
	fmt.Println(Variables())
	fmt.Print("\nЗадание 4\n")
	CalculateInt(4, 5)
	fmt.Print("\nЗадание 5\n")
	CalculateFloat(9.5, 5.21)
	fmt.Print("Задание 6\n")
	CalculateCount(4, 6, 9)
}
