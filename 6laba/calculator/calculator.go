package calculator

import (
	"fmt"
	"math/rand"
	"sync"
)

type calculator struct {
	A, B  float64
	Op    string
	ChRes chan float64
}

func calculate(a, b float64, op string) float64 {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			fmt.Print("Нельзя делить на ноль")
			return 0
		}
		return a / b
	default:
		fmt.Println("Неизвестная операция: ", op)
		return 0
	}
}

func Server(calculator <-chan calculator) {
	for c := range calculator {
		res := calculate(c.A, c.B, c.Op)
		c.ChRes <- res
	}
}

func Client(a, b float64, op string, calculators chan<- calculator, wg *sync.WaitGroup) {
	defer wg.Done()
	resCh := make(chan float64)
	calculators <- calculator{
		A:     a,
		B:     b,
		Op:    op,
		ChRes: resCh,
	}
	res := <-resCh
	fmt.Printf("%v %s %v = %.2v\n", a, op, b, res)
}

func Run() {
	var wg sync.WaitGroup
	c := make(chan calculator)
	go Server(c)
	ops := []string{"+", "-", "*", "/"}

	for i := 0; i < 10; i++ {
		op := ops[rand.Intn(len(ops))]
		a := float64(rand.Intn(20) + 1)
		b := float64(rand.Intn(20) + 1)
		wg.Add(1)
		go Client(a, b, op, c, &wg)

	}
	wg.Wait()
	close(c)
}
