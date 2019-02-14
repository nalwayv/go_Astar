package main

//---------------------------------------------------

import (
	"bufio"
	"os"
	"strings"

	"./gologger"
	"./grid"
	"./position"
)

// W / H
const (
	MapWidth  int = 10
	MapHeight int = 10
)

var (
	golog = gologger.GetInstance("src/gologger/golog.log")
)

func loadFromFile(fname string) []string {

	file, err := os.Open(fname)
	if err != nil {
		golog.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, MapWidth*MapHeight)
	idx := 0

	for scanner.Scan() {
		// break scanline up
		line := strings.Split(scanner.Text(), "")

		for _, v := range line {
			lines[idx] = v
			idx++
		}
	}

	return lines
}

func mainGrid() {

	levelMap := loadFromFile("assets/levelmap.txt")

	var g = grid.Grid{Width: MapWidth, Height: MapHeight}
	g.Cells = make([]grid.Cell, MapWidth*MapHeight)

	from := position.Pos{}
	to := position.Pos{}

	// setup
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {

			idx := g.GetCell(x, y)
			sym := string(levelMap[idx])

			g.Cells[idx] = grid.Cell{}
			g.Cells[idx].X = x
			g.Cells[idx].Y = y
			g.Cells[idx].Sym = sym

			// blocked
			if sym == "*" {
				g.Cells[idx].Blocked = true
			}

			// start/end pos
			if sym == "P" {
				from = g.Cells[idx].Pos
			}

			if sym == "E" {
				to = g.Cells[idx].Pos
			}

			// tile cost
			switch sym {
			case "~":
				g.Cells[idx].Cost = 50
			case "#":
				g.Cells[idx].Cost = 25
			default:
				g.Cells[idx].Cost = 1
			}
		}
	}

	// search
	path, found := g.Star(from, to)
	if found {
		//fmt.Println(path)

		for _, v := range path {
			from = from.Add(v.X, v.Y)
			idx := g.GetCell(from.X, from.Y)
			g.Cells[idx].Sym = string('.')
		}

	}

	// draw
	g.Draw()
}

//---------------------------------------------------

func main() {
	mainGrid()
}

//---------------------------------------------------
