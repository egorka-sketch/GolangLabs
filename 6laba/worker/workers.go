package worker

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Task struct {
	Id   int
	Line string
}
type Result struct {
	Id          int
	ReverseLine string
	WorkerId    int
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func worker(in <-chan Task, results chan<- Result, workerId int) {
	for task := range in {
		results <- Result{Id: task.Id, ReverseLine: Reverse(task.Line), WorkerId: workerId}
	}

}
func Run() {
	var wg sync.WaitGroup
	var lines []string
	var n int
	fmt.Print("Введите число воркеров: ")
	fmt.Scan(&n)
	if n < 1 {
		fmt.Print("Не может быть меньше 1")
		return
	}
	f, err := os.Open("6laba/input.txt")
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		lines = append(lines, s)
	}
	in := make(chan Task)
	go func() {
		defer close(in)
		for i, line := range lines {
			in <- Task{Id: i, Line: line}
		}
	}()
	out := make(chan Result)
	for w := 1; w <= n; w++ {
		wg.Go(func() {
			worker(in, out, w)
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	reversed := make([]string, len(lines))
	for res := range out {
		reversed[res.Id] = res.ReverseLine
		fmt.Printf("Воркер %d обработал строку %d : %s\n", res.WorkerId, res.Id+1, res.ReverseLine)

	}

}
