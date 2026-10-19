package channels

import (
	"fmt"
	"time"
)

func GenerateFibonacci(n int, ch chan<- int) {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a
		a, b = b, a+b
	}
	close(ch)
}

func PrintFibonacci(ch <-chan int, done chan<- bool) {
	for v := range ch {
		fmt.Println(v)
		time.Sleep(200 * time.Millisecond)
	}
	done <- true
}

func RunGoroutines() {
	fmt.Print("Задание 2\n")
	ch := make(chan int)
	done := make(chan bool)
	go GenerateFibonacci(10, ch)
	go PrintFibonacci(ch, done)
	<-done
}
