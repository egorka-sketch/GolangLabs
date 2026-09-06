package array

import "fmt"

func Slice() {
	arr := [5]int{10, 20, 30, 40, 50}
	fmt.Println("Массив из которого будем брать срез: ", arr)
	slice := arr[1:4]
	fmt.Println("Слайс: ", slice)
	slice = append(slice, 99, 11)
	fmt.Println("Слайс с новым элементом : ", slice)
	i := 2
	slice = append(slice[:i], slice[i+1:]...)
	fmt.Println("Слайс с удаленными элементом: ", slice)
}
