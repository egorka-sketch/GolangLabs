package person

import "fmt"

type Person struct {
	name string
	age  int
}

func NewPerson(name string, age int) *Person {
	return &Person{
		name,
		age}
}
func (p *Person) GetInfo() {
	fmt.Printf("%s, %d лет\n", p.name, p.age)
}
func (p *Person) BirthdayPerson() {
	p.age++
}
