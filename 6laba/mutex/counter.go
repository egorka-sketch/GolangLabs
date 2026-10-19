package mutex

import (
	"fmt"
	"sync"
)

func CounterOutMutex(wg *sync.WaitGroup, counter *int) {
	defer wg.Done()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		*counter++

	}
}

func CounterMutex(wg *sync.WaitGroup, mu *sync.Mutex, counter *int) {
	defer wg.Done()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		mu.Lock()
		*counter++
		mu.Unlock()
	}
}

func RunGoroutines() {
	var wg1, wg2 sync.WaitGroup
	var mu sync.Mutex
	counter1 := 0
	counter2 := 0
	wg1.Add(10)
	for i := 0; i < 10; i++ {
		go CounterOutMutex(&wg1, &counter1)
	}
	wg1.Wait()
	fmt.Print("Задание 4\n")
	fmt.Println("Без мьютекса из 10000 насчитано: ", counter1)
	wg2.Add(10)
	for i := 0; i < 10; i++ {
		go CounterMutex(&wg2, &mu, &counter2)
	}
	wg2.Wait()
	fmt.Println("С мьютексом из 10000 насчитано: ", counter2)
}
