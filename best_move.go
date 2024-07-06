package main

import (
	"math"
	"math/rand"
)

//////// PLAYER2 is the AI ///////////

type MoveFinder struct{}

// Utility function to copy the board
func copyBoard(board [][]int32) [][]int32 {
	copy := make([][]int32, ROWS)
	for row := range board {
		copy[row] = make([]int32, COLS)
	}

	for row := range board {
		for col := range board[0] {
			copy[row][col] = board[row][col]
		}
	}
	return copy
}

// Function to score the clump
func (f *MoveFinder) scoreClump(clump []int32) float64 {
	var countPlayer1 int32 = f.countClump(clump, PLAYER1)
	var countPlayer2 int32 = f.countClump(clump, PLAYER2)
	var countEmpty int32 = f.countClump(clump, EMPTY)

	if countPlayer1 == 3 && countEmpty == 1 {
		return -10.
	} else if countPlayer2 == 3 && countEmpty == 1 {
		return 10.
	} else if countPlayer1 == 2 && countEmpty == 2 {
		return -5.
	} else if countPlayer2 == 2 && countEmpty == 2 {
		return 5.
	} else if countPlayer1 == 1 && countEmpty == 3 {
		return -1.
	} else if countPlayer2 == 1 && countEmpty == 3 {
		return 1.
	}

	return 0.
}

// Function to count the number of occurances in a clump
func (f *MoveFinder) countClump(clump []int32, value int32) int32 {
	var count int32 = 0
	for i := range clump {
		if clump[i] == value {
			count++
		}
	}
	return count
}

// Writing the score funciton
func (f *MoveFinder) score(board [][]int32) float64 {
	// If the game is a tie: return 0: Check if the game is a tie
	if len(f.getValidCols(board)) == 0 {
		return 0.
	}
	// If the game is a win for player1: return -100_000_000 else return +100_000_000
	if f.checkWinner(board) == PLAYER1 {
		return -100_000_000
	} else if f.checkWinner(board) == PLAYER2 {
		return +100_000_000
	}

	// Score for other positions:
	clump := make([]int32, 4)
	var score float64 = 0

	// Check the horizontal lines
	for row := range board {
		for col := 0; col < COLS-3; col++ {
			// Populate the clump array
			for i := range clump {
				clump[i] = board[row][col+i]
			}
			score += f.scoreClump(clump)
		}
	}

	// Check Vertical lines
	for col := range board[0] {
		for row := 0; row < ROWS-3; row++ {
			for i := range clump {
				clump[i] = board[row+i][col]
			}
			score += f.scoreClump(clump)
		}
	}

	// check negative slope
	for row := 0; row < ROWS-3; row++ {
		for col := 0; col < COLS-3; col++ {
			for i := range clump {
				clump[i] = board[row+i][col+i]
			}
			score += f.scoreClump(clump)
		}
	}

	// Check positive slope
	for row := ROWS - 1; row > ROWS-4; row-- {
		for col := 0; col < COLS-3; col++ {
			for i := range clump {
				clump[i] = board[row-i][col+i]
			}
			score += f.scoreClump(clump)
		}
	}
	return score
}

func (f *MoveFinder) BestMove(originalBoard [][]int32, depth int32) int32 {
	// board := copyBoard(originalBoard)
	// validCols := f.getValidCols(board)

	// score := math.Inf(-1)
	// bestColumn := validCols[0]
	// for _, col := range validCols {
	// 	// Copy the original board
	// 	anotherBoard := copyBoard(originalBoard)
	// 	// Get the valid row for the column
	// 	row := f.getRow(anotherBoard, col)
	// 	// Make the move
	// 	anotherBoard[row][col] = PLAYER2
	// 	// Score the move
	// 	colScore := f.score(anotherBoard)
	// 	// Update teh score
	// 	if colScore > score {
	// 		score = colScore
	// 		bestColumn = col
	// 	}
	// }

	// return bestColumn
	col, _ := f.Minimax(originalBoard, depth, true)
	return col
}

// Function to get the row
func (f *MoveFinder) getRow(board [][]int32, col int32) int32 {
	for row := len(board) - 1; row >= 0; row-- {
		if board[row][col] == EMPTY {
			return int32(row)
		}
	}
	return -1
}

func (f *MoveFinder) getValidCols(board [][]int32) []int32 {
	var validCols []int32
	for col := range board[0] {
		if f.isColumnValid(int32(col), board) {
			validCols = append(validCols, int32(col))
		}
	}
	return validCols
}

func (f *MoveFinder) isColumnValid(col int32, board [][]int32) bool {
	return board[0][col] == EMPTY
}

// Funciton to check if there is a winner
func (f *MoveFinder) checkWinner(board [][]int32) int32 {
	// check horizontal
	for row := range board {
		for col := 0; col < COLS-3; col++ {
			if board[row][col] != EMPTY && board[row][col] == board[row][col+1] && board[row][col+1] == board[row][col+2] && board[row][col+2] == board[row][col+3] {
				return board[row][col]
			}
		}
	}
	// check vertical
	for col := range board[0] {
		for row := 0; row < ROWS-3; row++ {
			if board[row][col] != EMPTY && board[row][col] == board[row+1][col] && board[row+1][col] == board[row+2][col] && board[row+2][col] == board[row+3][col] {
				return board[row][col]
			}
		}
	}

	// check negative slope
	for row := 0; row < ROWS-3; row++ {
		for col := 0; col < COLS-3; col++ {
			if board[row][col] != EMPTY && board[row][col] == board[row+1][col+1] && board[row+1][col+1] == board[row+2][col+2] && board[row+2][col+2] == board[row+3][col+3] {
				return board[row][col]
			}
		}
	}

	// check positive slope
	for row := ROWS - 1; row > ROWS-4; row-- {
		for col := 0; col < COLS-3; col++ {
			if board[row][col] != EMPTY && board[row][col] == board[row-1][col+1] && board[row-1][col+1] == board[row-2][col+2] && board[row-2][col+2] == board[row-3][col+3] {
				return board[row][col]
			}
		}
	}

	return -1
}

// Writing the minimax function
func (f *MoveFinder) Minimax(originalBoard [][]int32, depth int32, maximizingPlayer bool) (int32, float64) {
	validCols := f.getValidCols(originalBoard)
	winner := f.checkWinner(originalBoard)

	// Check if the match is a tie or has a winner : This is when the recursion reaches a leaf node
	if len(validCols) <= 0 || winner != -1 || depth <= 0 {
		return -1, f.score(originalBoard)
	}

	// Write the recursive function
	var board [][]int32
	var bestColumn int32
	var bestScore float64

	if maximizingPlayer {
		// Initialize bestscore
		bestScore = math.Inf(-1)
		bestColumn = validCols[rand.Int31n(int32(len(validCols)))]

		// Loop through all the valid columns
		for i := range validCols {
			// Copy the board
			board = copyBoard(originalBoard)
			// Get the row from the column
			col := validCols[i]
			row := f.getRow(board, col)
			// Make the move
			board[row][col] = PLAYER2
			_, score := f.Minimax(board, depth-1, !maximizingPlayer)

			if score > bestScore {
				bestScore = score
				bestColumn = col
			}
		}
	} else {
		// Initialize bestscore
		bestScore = math.Inf(1)
		bestColumn = validCols[rand.Int31n(int32(len(validCols)))]

		// Loop through all the valid columns
		for i := range validCols {
			// Copy the board
			board = copyBoard(originalBoard)
			// Get the row from the column
			col := validCols[i]
			row := f.getRow(board, col)
			// Make the move
			board[row][col] = PLAYER1
			_, score := f.Minimax(board, depth-1, !maximizingPlayer)

			if score < bestScore {
				bestScore = score
				bestColumn = col
			}
		}
	}

	return bestColumn, bestScore
}
