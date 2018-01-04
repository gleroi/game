package main

import (
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

type Direction int

const (
	DirR Direction = 0
	DirD Direction = 1
	DirL Direction = 2
	DirU Direction = 3
)

type player struct {
	picture   pixel.Picture
	anim      int
	direction Direction
	position  pixel.Vec
	sprite    *pixel.Sprite
}

func (p *player) Draw(target pixel.Target) {
	const spriteCols = 9
	const spriteRows = 4
	spriteWidth := p.picture.Bounds().Max.X / spriteCols
	spriteHeight := p.picture.Bounds().Max.Y / spriteRows
	spriteX := float64(1+p.anim%(spriteCols-1)) * spriteWidth
	spriteY := float64(p.anim/spriteCols) * spriteHeight

	bounds := pixel.R(spriteX, spriteY, spriteX+spriteWidth, spriteY+spriteHeight)
	if p.sprite == nil {
		p.sprite = pixel.NewSprite(p.picture, bounds)
	} else {
		p.sprite.Set(p.picture, bounds)
	}
	p.sprite.Draw(target, pixel.IM.Moved(p.position))
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
		picture:   pic,
		anim:      0,
		direction: DirR,
		position:  win.Bounds().Center(),
	}

	last := time.Now()
	for !win.Closed() {
		now := time.Now()
		elapsed := now.Sub(last)

		if win.Pressed(pixelgl.KeyLeft) {
			p.direction = DirL
		}
		if win.Pressed(pixelgl.KeyRight) {
			p.direction = DirR
		}
		if win.Pressed(pixelgl.KeyDown) {
			p.direction = DirD
		}
		if win.Pressed(pixelgl.KeyUp) {
			p.direction = DirU
		}

		if elapsed >= 55*time.Millisecond {
			p.anim = int(p.direction)*9 + (p.anim+1)%9
			last = now

		}

		win.Clear(colornames.Burlywood)
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
