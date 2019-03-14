package main

import (
	"image"
	"log"
	"math/cmplx"

	"github.com/hajimehoshi/ebiten"
)

//World is a struct that holds a 2-D string of bools, turning a pixel on or off, a width, and a height
type World struct {
	area   [][]int
	width  int
	height int
}

const (
	screenWidth     = 500
	screenHeight    = 500
	scale           = 2
	titleString     = "Go Frac!"
	MAX_ITTERATIONS = 50
)

var (
	world       *World
	imageBuffer *image.RGBA
	zoom        int
	focus       complex128
)

func update(screen *ebiten.Image) error {
	// When IsDrawingSkipped is true, the rendered result is not adopted.
	// Skip rendering then.
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	//button scan
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		zoom++
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		zoom--
	}
	//screen update
	world.Progress()
	world.DrawImage(imageBuffer)
	screen.ReplacePixels(imageBuffer.Pix)
	return nil
}

func main() {
	//init
	zoom = 1.0
	focus = complex(0, 0)
	world = NewWorld(screenWidth, screenHeight)
	imageBuffer = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, scale, titleString); err != nil {
		log.Fatal(err)
	}
}

// NewWorld creates a new world screen. This is the buffer for the output window. Update it between ticks.
func NewWorld(width, height int) *World {
	world := World{}
	world.area = makeArea(width, height)
	world.width = width
	world.height = height
	return &world
}

// makeArea initializes the matrix of pixels for the world.
func makeArea(width, height int) [][]int {
	area := make([][]int, height)
	for i := 0; i < height; i++ {
		area[i] = make([]int, width)
	}
	return area
}

// DrawImage paints current game state
func (w *World) DrawImage(img *image.RGBA) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			w.drawPixel(img, x, y)
		}
	}
}

func (w *World) drawPixel(img *image.RGBA, x, y int) {
	pos := 4*y*w.width + 4*x
	if w.area[y][x] == 0 {
		img.Pix[pos] = 0
		img.Pix[pos+1] = 0
		img.Pix[pos+2] = 0
		img.Pix[pos+3] = 0
		return
	} else if w.area[y][x] == 2 {
		img.Pix[pos] = 0xff
		img.Pix[pos+1] = 0
		img.Pix[pos+2] = 0
		img.Pix[pos+3] = 0xff
		return
	} else if w.area[y][x] == 3 {
		img.Pix[pos] = 0
		img.Pix[pos+1] = 0xff
		img.Pix[pos+2] = 0
		img.Pix[pos+3] = 0xff
		return
	} else if w.area[y][x] == 4 {
		img.Pix[pos] = 0
		img.Pix[pos+1] = 0
		img.Pix[pos+2] = 0xff
		img.Pix[pos+3] = 0xff
		return
	} else if w.area[y][x] == 5 {
		img.Pix[pos] = 0xff
		img.Pix[pos+1] = 0
		img.Pix[pos+2] = 0xff
		img.Pix[pos+3] = 0xff
		return
	} else {
		img.Pix[pos] = 0xff
		img.Pix[pos+1] = 0xff
		img.Pix[pos+2] = 0
		img.Pix[pos+3] = 0xff
	}
}

// Progress game state by one tick
func (w *World) Progress() {
	next := makeArea(w.width, w.height)

	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			next[y][x] = fractalValue(normalizeCoords(x, y))
		}
	}
	w.area = next
}

//noramalizeCoords takes an absolute pixel position, zoom, and focus, and returns a complex number for that pixel. This should then be run through fractalValue.
func normalizeCoords(x, y int) complex128 {
	real := (float64(2*x)/float64(zoom*screenWidth) - 1.0)
	im := (float64(2*y)/float64(zoom*screenHeight) - 1.0)
	return complex(real, im) - focus
}

//fractalValue takes a complex number and checks to see if it diverges from the Julia set of c in MAX_ITERATIONS.
func fractalValue(c complex128) int {
	z := complex(0.0, 0.0)
	for i := 0; i < MAX_ITTERATIONS; i++ {
		if cmplx.Abs(z) > 2.0 {
			return i
		}
		z = z*z + c
	}
	return 0
}
