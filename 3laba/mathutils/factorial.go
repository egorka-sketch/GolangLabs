package mathutils

func Factorial(n int) int {
	if n <= 0 {
		return 1
	}
	res := 1
	for i := 2; i <= n; i++ {
		res = res * i
	}
	return res
}
