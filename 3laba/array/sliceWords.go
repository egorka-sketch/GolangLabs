package array

import (
	"fmt"
	"unicode/utf8"
)

func SliceWord() {
	words := []string{"Golang", "Rust", "Java", "Ruby",
		"Node", "C#"}
	longestWord := ""
	for _, word := range words {
		if utf8.RuneCountInString(word) > utf8.RuneCountInString(longestWord) {
			longestWord = word
		}
	}
	fmt.Println(longestWord)
}
