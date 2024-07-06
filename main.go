package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	game := Game{}
	game.init()
	defer game.close()

	game.loop()
}

var moveFinder MoveFinder
var myMove MoveFinder

type Game struct {
	running bool
	over    bool
	board   [][]int32
	Player1 Player
	Player2 Player
	Players []Player
	Turn    int32
	winner  int32
}

func (g *Game) init() {
	g.running = true
	g.over = false
	g.winner = EMPTY

	rl.InitWindow(WIDTH, HEIGHT, TITLE)
	rl.SetTargetFPS(FPS)

	g.board = make([][]int32, ROWS)
	for row := range g.board {
		g.board[row] = make([]int32, COLS)
	}

	g.Player1 = Player{Id: PLAYER1, Color: rl.Red, VictoryPhrase: "Player 1 Wins"}
	g.Player2 = Player{Id: PLAYER2, Color: rl.Yellow, VictoryPhrase: "Player 2 Wins"}
	g.Players = []Player{g.Player1, g.Player2}

	g.Turn = g.Player1.Id

}

func (g *Game) close() {
	rl.CloseWindow()
}

func (g *Game) loop() {
	for g.running {
		if rl.WindowShouldClose() {
			g.running = false
			return
		}
		g.update()
		g.draw()
	}
}

func (g *Game) reset() {
	// g.over = false
	// g.running = true
	// g.Turn = g.Player1.Id
	// g.winner = EMPTY

	// for row := range g.board {
	// 	for col := range g.board[row] {
	// 		g.board[row][col] = EMPTY
	// 	}
	// }
	g.close()
	g.init()
}

func (g *Game) makeMove(row, col int32) {
	g.dropPiece(row, col)

	if g.checkWinner() {
		g.over = true
		g.winner = g.Turn
	}

	// Check game over with tie
	if len(g.getValidCols()) == 0 {
		g.over = true
		g.winner = EMPTY
	}

	// switch column
	if g.Turn == g.Player1.Id {
		g.Turn = g.Player2.Id
	} else {
		g.Turn = g.Player1.Id
	}
}

func (g *Game) update() {
	if g.over {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			g.reset()
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && g.Turn == g.Player1.Id {
		// Drop down effect
		col := g.GetCol()
		if g.isColumnValid(col) {

			row := g.getRow(col)
			g.makeMove(row, col)
		}
	}
	// if g.Turn == g.Player1.Id && !g.over {
	// 	col := myMove.BestMove(g.board, 4)
	// 	row := g.getRow(col)
	// 	g.makeMove(row, col)
	// }

	if g.Turn == g.Player2.Id && !g.over {
		// AI Move
		col := moveFinder.BestMove(g.board, 5)
		row := g.getRow(col)
		g.makeMove(row, col)
	}
}

func (g *Game) draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(BgColor)

	rl.DrawRectangle(0, GridSize, WIDTH, HEIGHT-GridSize, rl.Blue)

	for row := range g.board {
		for col := range g.board[row] {
			// Draw the circles in the board
			position := g.board[row][col]
			if position == EMPTY {
				rl.DrawCircle(int32(GridSize*col+GridSize/2), int32((GridSize*row+GridSize/2)+GridSize), RADIUS, rl.Black)
			} else {
				rl.DrawCircle(int32(GridSize*col+GridSize/2), int32((GridSize*row+GridSize/2)+GridSize), RADIUS, g.Players[position-1].Color)
			}
		}
	}

	rl.DrawCircle(rl.GetMouseX(), GridSize/2, RADIUS, g.Players[g.Turn-1].Color)

	if g.over {
		if g.winner == EMPTY {
			rl.DrawText("Tie!", 20, 20, 40, rl.White)
		} else {
			rl.DrawText(g.Players[g.winner-1].VictoryPhrase, 20, 20, 40, g.Players[g.winner-1].Color)
		}
	}
}

func (g *Game) GetCol() int32 {
	return int32(rl.GetMouseX() / 100)
}

func (g *Game) isColumnValid(col int32) bool {
	return g.board[0][col] == EMPTY
}

func (g *Game) getRow(col int32) int32 {
	for row := len(g.board) - 1; row >= 0; row-- {
		if g.board[row][col] == EMPTY {
			return int32(row)
		}
	}
	return -1
}

func (g *Game) dropPiece(row, col int32) {
	g.board[row][col] = g.Turn
}

func (g *Game) getValidCols() []int32 {
	var validCols []int32
	for col := range g.board[0] {
		if g.isColumnValid(int32(col)) {
			validCols = append(validCols, int32(col))
		}
	}
	return validCols
}

func (g *Game) checkWinner() bool {
	// check horizontal
	for row := range g.board {
		for col := 0; col < COLS-3; col++ {
			if g.board[row][col] != EMPTY && g.board[row][col] == g.board[row][col+1] && g.board[row][col+1] == g.board[row][col+2] && g.board[row][col+2] == g.board[row][col+3] {
				return true
			}
		}
	}
	// check vertical
	for col := range g.board[0] {
		for row := 0; row < ROWS-3; row++ {
			if g.board[row][col] != EMPTY && g.board[row][col] == g.board[row+1][col] && g.board[row+1][col] == g.board[row+2][col] && g.board[row+2][col] == g.board[row+3][col] {
				return true
			}
		}
	}

	// check negative slope
	for row := 0; row < ROWS-3; row++ {
		for col := 0; col < COLS-3; col++ {
			if g.board[row][col] != EMPTY && g.board[row][col] == g.board[row+1][col+1] && g.board[row+1][col+1] == g.board[row+2][col+2] && g.board[row+2][col+2] == g.board[row+3][col+3] {
				return true
			}
		}
	}

	// check positive slope
	for row := ROWS - 1; row > ROWS-4; row-- {
		for col := 0; col < COLS-3; col++ {
			if g.board[row][col] != EMPTY && g.board[row][col] == g.board[row-1][col+1] && g.board[row-1][col+1] == g.board[row-2][col+2] && g.board[row-2][col+2] == g.board[row-3][col+3] {
				return true
			}
		}
	}

	return false
}
