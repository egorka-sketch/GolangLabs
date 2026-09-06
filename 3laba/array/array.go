package array

import "math/rand/v2"

func GenerateArray() [5]int {
	var arr [5]int
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.IntN(100)
	}
	return arr
}
