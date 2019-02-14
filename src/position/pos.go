package position

import (
	"errors"
	"math"
)

// Pos ...
type Pos struct {
	X int
	Y int
}

// Add ...
func (p Pos) Add(x, y int) Pos {
	return Pos{X: p.X + x, Y: p.Y + y}
}

// Sub ...
func (p Pos) Sub(x, y int) Pos {
	return Pos{X: p.X - x, Y: p.Y - y}
}

// Scale ...
func (p Pos) Scale(val int) Pos {
	return Pos{X: p.X * val, Y: p.Y * val}
}

// Div ...
func (p Pos) Div(val int) (Pos, error) {
	if val == 0 {
		return Pos{}, errors.New("devide by zero")
	}

	return Pos{X: p.X / val, Y: p.Y / val}, nil
}

// Length ...
func (p Pos) Length() int {
	return int(math.Sqrt(float64(p.X*p.X) + float64(p.Y*p.Y)))
}
