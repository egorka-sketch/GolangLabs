package main

import (
	"5laba/book"
	"5laba/geometry"
	"5laba/person"
	"fmt"
)

func main() {
	p := person.NewPerson("Алексей", 55)
	p.GetInfo()
	p.BirthdayPerson()
	p.GetInfo()
	c := geometry.NewCircle(5)
	fmt.Printf("Площадь круга = %.2f\n", c.Area())
	r := geometry.NewRectangle(4.33, 5)
	fmt.Println("Площадь треугольника =  ", r.Area())

	shapes := []geometry.Shape{
		geometry.NewCircle(5.0),
		geometry.NewRectangle(4.0, 6.0),
		geometry.NewCircle(2.5),
	}

	geometry.PrintAreas(shapes)

	b := book.NewBook("Безымянка", "Сергей Палий", 2013)
	fmt.Println(*b)
}
