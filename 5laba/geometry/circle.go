package geometry

import (
	"fmt"
	"math"
)

type Circle struct {
	radius float64
}
type Rectangle struct {
	width  float64
	height float64
}
type Shape interface {
	Area() float64
}

func NewCircle(radius float64) *Circle {
	return &Circle{radius}
}
func NewRectangle(width, height float64) *Rectangle {
	return &Rectangle{
		width,
		height}
}

func (c *Circle) Area() float64 {
	return math.Pi * math.Pow(c.radius, 2)
}
func (r *Rectangle) Area() float64 {
	return r.width * r.height
}

func PrintAreas(shapes []Shape) {
	for i, s := range shapes {
		fmt.Println("Фигура", i+1, s.Area())
	}
}
