package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	WIDTH  = 700
	HEIGHT = 700

	TITLE = "Connect 4"

	GridSize = 100
	ROWS     = 6
	COLS     = 7

	EMPTY   = 0
	FPS     = 30
	PLAYER1 = 1
	PLAYER2 = 2

	RADIUS = GridSize * 5 / 11
)

var BgColor rl.Color = rl.Black
