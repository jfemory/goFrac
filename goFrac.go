package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
	// When IsDrawingSkipped is true, the rendered result is not adopted.
	// Skip rendering then.
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	//Update state.
	ebitenutil.DebugPrint(screen, "fractal goes here")
	return nil
}

func main() {
	ebiten.Run(update, 320, 240, 2, "Go Frac!")
}
