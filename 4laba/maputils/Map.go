package maputils

import "fmt"

func CreateMap() map[string]int {
	m := map[string]int{
		"Оля":    23,
		"Альпак": 52,
	}

	m["Максим"] = 71
	return m
}

func AverageYears(m map[string]int) float64 {
	sum := 0
	for _, v := range m {
		sum += v
	}
	return float64(sum) / float64(len(m))
}

func DeletePerson(m map[string]int, name string) {
	fmt.Println("Задание 3\nдо удаления :", m)
	delete(m, name)
	fmt.Println("После удаления: ", m)
}
