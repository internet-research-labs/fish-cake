package server

import (
	"fmt"
	"strings"
)

type GameOfLife struct {
	Cells [][]int
	NRows int
	NCols int
}

func Board(rows, cols int) [][]int {
	cells := make([][]int, rows)
	for i := int(0); i < rows; i++ {
		cells[i] = make([]int, cols)
	}
	return cells
}

func NewGameOfLife(rows, cols int) *GameOfLife {

	return &GameOfLife{
		Cells: Board(rows, cols),
		NRows: rows,
		NCols: cols,
	}
}

func (game *GameOfLife) InBounds(i, j int) bool {
	return i >= 0 && i < game.NRows && j >= 0 && j < game.NCols
}

func (game *GameOfLife) LivingCount(y, x int) int {

	count := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {

			u := y + i
			v := x + j

			if !game.InBounds(u, v) {
				continue
			}

			if i == 0 && j == 0 {
				continue
			}

			count += game.Cells[u][v]
		}
	}

	return count
}

func (game *GameOfLife) Add(i, j int) error {
	if !game.InBounds(i, j) {
		return fmt.Errorf("Out of bounds")
	}
	game.Cells[i][j] = 1
	return nil
}

func (game *GameOfLife) Tick() {
	board := Board(game.NRows, game.NCols)

	for i := 0; i < game.NRows; i++ {
		for j := 0; j < game.NCols; j++ {

			count := game.LivingCount(i, j)

			var val int
			alive := game.Cells[i][j] == 1

			if alive {
				if count == 2 || count == 3 {
					val = 1
				}
			} else {
				if count == 3 {
					val = 1
				}
			}

			board[i][j] = val
		}
	}

	game.Cells = board
}

func (game *GameOfLife) String() string {

	var builder strings.Builder

	builder.WriteString("+")
	builder.WriteString(strings.Repeat("-", int(game.NCols)))
	builder.WriteString("+\n")

	for _, row := range game.Cells {
		_ = row
		builder.WriteString("|")
		for _, r := range row {
			var c string
			if r == 1 {
				c = "*"
			} else {
				c = " "
			}
			builder.WriteString(fmt.Sprintf("%v", c))
		}
		builder.WriteString("|\n")
	}

	builder.WriteString("+")
	builder.WriteString(strings.Repeat("-", int(game.NCols)))
	builder.WriteString("+")

	return builder.String()
}
