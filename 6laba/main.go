package main

import (
	"6laba/calculator"
	"6laba/channels"
	"6laba/goroutines"
	"6laba/mutex"
	"6laba/selector"
	"fmt"
)

func main() {
	goroutines.Goroutines()
	channels.RunGoroutines()
	selector.Selector()
	mutex.RunGoroutines()
	fmt.Print("\nЗадание 5\n")
	calculator.Run()
}
