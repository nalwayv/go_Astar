package grid

import (
	"errors"
	"fmt"
	"math"

	"../gologger"
	"../position"
	"../pq"
)

//---------------------------------------------------
// MAP INFO
// 'S' -> start
// 'E' -> end
// '.' -> path
// ' ' -> open
// '*' -> blocked
// '~' -> water
// '#' -> sand
//---------------------------------------------------

// W / H ...

var (
	golog = gologger.GetInstance("src/gologger/golog.log")
)

//---------------------------------------------------

// Cell ...
type Cell struct {
	position.Pos
	Sym     string
	Visited bool
	Blocked bool
	Parent  *Cell
	Cost    int
}

//---------------------------------------------------

// Grid ...
type Grid struct {
	Cells  []Cell
	Width  int
	Height int
}

// GetCell ...
func (g Grid) GetCell(x, y int) int {
	return y*(g.Width) + x
}

func (g Grid) canMove(pos position.Pos) (bool, error) {
	idx := g.GetCell(pos.X, pos.Y)

	if !g.Cells[idx].Blocked && pos.X >= 0 && pos.X < g.Width && pos.Y >= 0 && pos.Y < g.Height {
		return true, nil
	}

	return false, errors.New("out of range")
}

func (g Grid) getNeighbours(pos position.Pos) []position.Pos {

	neighbours := make([]position.Pos, 0, 4)

	up := position.Pos{X: pos.X, Y: pos.Y - 1}
	down := position.Pos{X: pos.X, Y: pos.Y + 1}
	left := position.Pos{X: pos.X - 1, Y: pos.Y}
	right := position.Pos{X: pos.X + 1, Y: pos.Y}

	upRight := position.Pos{X: pos.X + 1, Y: pos.Y - 1}
	upLeft := position.Pos{X: pos.X - 1, Y: pos.Y - 1}
	downRight := position.Pos{X: pos.X + 1, Y: pos.Y + 1}
	downLeft := position.Pos{X: pos.X - 1, Y: pos.Y + 1}

	if ok, _ := g.canMove(up); ok {
		neighbours = append(neighbours, up)
	}

	if ok, _ := g.canMove(down); ok {
		neighbours = append(neighbours, down)
	}

	if ok, _ := g.canMove(left); ok {
		neighbours = append(neighbours, left)
	}

	if ok, _ := g.canMove(right); ok {
		neighbours = append(neighbours, right)
	}

	// diag
	if ok, _ := g.canMove(upRight); ok {
		neighbours = append(neighbours, upRight)
	}

	if ok, _ := g.canMove(upLeft); ok {
		neighbours = append(neighbours, upLeft)
	}

	if ok, _ := g.canMove(downRight); ok {
		neighbours = append(neighbours, downRight)
	}

	if ok, _ := g.canMove(downLeft); ok {
		neighbours = append(neighbours, downLeft)
	}

	return neighbours
}

// func (g *grid) bfs(from, to position.Pos) {
// 	frontier := make([]position.Pos, 0, 8)
// 	frontier = append(frontier, from)
// 	visited := make(map[position.Pos]bool)
// 	visited[from] = true
// 	found := false

// 	for len(frontier) > 0 {
// 		current := frontier[0]
// 		frontier = frontier[1:]

// 		if current == to {
// 			found = true
// 			break
// 		}

// 		for _, next := range g.getNeighbours(current) {
// 			if !visited[next] {
// 				frontier = append(frontier, next)
// 				visited[next] = true

// 			}
// 		}
// 	}

// 	if found {
// 		fmt.Println("found")
// 		for k := range visited {
// 			idx := g.getCell(k.X, k.Y)
// 			g.cells[idx].sym = "."
// 		}
// 	} else {
// 		fmt.Println("not found")
// 	}
// }

// Star ...
func (g *Grid) Star(from, to position.Pos) ([]position.Pos, bool) {
	// check if from and to are ok to use
	// ok, err := g.canMove(from)
	dirs := make([]position.Pos, 0)

	if ok, err := g.canMove(from); !ok && err != nil {
		golog.Fatalln("'from' not provided for star()")
		return dirs, false
	}

	if ok, err := g.canMove(to); !ok && err != nil {
		golog.Fatalln("'to' not provided for star()")
		return dirs, false
	}

	// manhattan heusistic
	heur := func(a, b position.Pos) int {
		dx := math.Abs(float64(a.X - b.X))
		dy := math.Abs(float64(a.Y - b.Y))
		return int(dx + dy)
	}

	//found := false

	frontier := make(pq.PQueue, 0, 8)
	frontier = frontier.Push(from, 1)

	// path back
	cameFrom := make(map[position.Pos]position.Pos)
	cameFrom[from] = from

	// movement cost
	costSoFar := make(map[position.Pos]int)
	costSoFar[from] = 0

	// search
	var current position.Pos
	for frontier.Len() > 0 {

		frontier, current = frontier.Pop()

		// found goal / path back
		if current == to {
			golog.Println("found a path")
			p := cameFrom[current]

			for p != from {

				dirs = append(dirs, position.Pos{X: p.X - cameFrom[p].X, Y: p.Y - cameFrom[p].Y})
				p = cameFrom[p]
			}

			// reverse list
			start := 0
			end := len(dirs) - 1
			for start < end {
				dirs[start], dirs[end] = dirs[end], dirs[start]
				start++
				end--
			}

			return dirs, true
		}

		// check neighbours
		for _, next := range g.getNeighbours(current) {

			// cost to travel plus tile cost
			idx := g.GetCell(current.X, current.Y)
			newCost := costSoFar[current] + g.Cells[idx].Cost

			// check is within cost so far
			_, ok := costSoFar[next]

			// if not in cost so far or cost is lower
			if !ok || newCost < costSoFar[next] {
				costSoFar[next] = newCost

				newPriority := newCost + heur(to, next)

				frontier = frontier.Push(next, newPriority)

				cameFrom[next] = current
			}
		}
	}

	golog.Println("did not find a path")
	return dirs, false
}

// Draw ...
func (g *Grid) Draw() {
	// draw map
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {

			idx := g.GetCell(x, y)
			sym := g.Cells[idx].Sym
			fmt.Print(sym, " ")

		}

		fmt.Println("")
	}
}
