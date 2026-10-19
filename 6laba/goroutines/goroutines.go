package goroutines

import (
	"fmt"
	"math/rand"
	"time"
)

func Factorial(n int) {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
		time.Sleep(time.Millisecond * 150)
	}
	fmt.Printf("\nФакториал от %d равен %d\n", n, result)
}
func Random(count int) {
	for i := 0; i < count; i++ {
		num := rand.Intn(100) + 1
		fmt.Println("\nСлучайное число: ", num)
		time.Sleep(time.Millisecond * 200)
	}
}

func Sum(n int) {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println("\nСумма ряда: ", sum)
}

func Goroutines() {
	fmt.Print("Задание 1")
	go Factorial(5)
	go Random(4)
	go Sum(10)

	time.Sleep(2 * time.Second)
}
