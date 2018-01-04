package main

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

type player struct {
	picture  pixel.Picture
	anim     int
	position pixel.Vec
}

func (p *player) Draw(target pixel.Target) {
	const spriteCols = 9
	const spriteRows = 4
	spriteWidth := p.picture.Bounds().Max.X / spriteCols
	spriteHeight := p.picture.Bounds().Max.Y / spriteRows
	spriteX := float64(p.anim%spriteCols) * spriteWidth
	spriteY := float64(p.anim/spriteCols) * spriteHeight

	bounds := pixel.R(spriteX, spriteY, spriteWidth, spriteHeight)
	fmt.Printf("i: %d, offset: %v\n", p.anim, bounds)
	sprite := pixel.NewSprite(p.picture, bounds)
	sprite.Draw(target, pixel.IM.Moved(p.position))
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Pixel Rocks!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	pic, err := loadPicture("professor_walk.png")
	if err != nil {
		panic(err)
	}

	p := player{
		picture:  pic,
		anim:     0,
		position: win.Bounds().Center(),
	}

	last := time.Now()
	for !win.Closed() {

		now := time.Now()
		elapsed := now.Sub(last)

		if elapsed >= 1000*time.Millisecond {
			p.anim = (p.anim + 1) % (4 * 9)
			last = now

		}

		win.Clear(colornames.Greenyellow)
		p.Draw(win)
		win.Update()
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
