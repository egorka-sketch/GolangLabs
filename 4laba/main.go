package main

import (
	"4laba/io"
	"4laba/maputils"
	"fmt"
)

func main() {
	m := maputils.CreateMap()
	fmt.Print("Задание 1\nСписок людей\n")
	for k, v := range m {
		fmt.Printf("%s %d года\n", k, v)
	}
	avg := maputils.AverageYears(m)
	fmt.Printf("Задание 2\nСреднее значение = %.2f\n", avg)

	maputils.DeletePerson(m, "Максим")

	io.Uppercase()
	io.SumSlice()
	io.InvertedSlice()
}
