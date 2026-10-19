package array

import "fmt"

func Slice() {
	arr := [6]int{10, 20, 30, 40, 50, 222}
	fmt.Println("Массив из которого будем брать срез: ", arr)
	slice := arr[1:5]
	fmt.Println("Слайс: ", slice)
	slice = append(slice, 99, 11)
	fmt.Println("Слайс с новым элементом : ", slice)
	i := 3
	slice = append(slice[:i], slice[i+1:]...)
	fmt.Println("Слайс с удаленными элементом: ", slice)
}
