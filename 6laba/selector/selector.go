package selector

import (
	"fmt"
	"math/rand"
	"time"
)

func GeneratorEven(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- (rand.Intn(50) + 1) * 2
		time.Sleep(time.Millisecond * 200)
	}
}

func GeneratorOdd(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- (rand.Intn(50))*2 + 1
		time.Sleep(time.Millisecond * 200)
	}
}

func Selector() {
	fmt.Print("\nЗадание 3\n")
	evenCh := make(chan int)
	oddCh := make(chan int)

	go GeneratorEven(evenCh)
	go GeneratorOdd(oddCh)

	for i := 0; i < 20; i++ {
		select {
		case e := <-evenCh:
			fmt.Println("Число четное: ", e)
		case o := <-oddCh:
			fmt.Println("Число нечетное: ", o)
		}
	}

}
