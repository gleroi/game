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

type Direction int

const (
	DirR Direction = iota
	DirD
	DirL
	DirU
)

type player struct {
	picture   pixel.Picture
	anim      int
	animStep  float64
	direction Direction
	position  pixel.Vec
	sprite    *pixel.Sprite
}

const spriteCols = 9
const spriteRows = 4

func (p *player) Draw(target pixel.Target) {
	spriteWidth := p.picture.Bounds().Max.X / spriteCols
	spriteHeight := p.picture.Bounds().Max.Y / spriteRows

	animOffset := int(p.direction)*spriteCols + p.anim
	spriteX := float64(animOffset%spriteCols) * spriteWidth
	spriteY := float64(animOffset/spriteCols) * spriteHeight

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
	animInterval := 4.5
	speed := 76.0
	for !win.Closed() {
		now := time.Now()
		elapsed := now.Sub(last)

		deltaM := 1.0
		var move pixel.Vec
		if win.Pressed(pixelgl.KeyLeft) {
			p.direction = DirL
			move = pixel.V(-deltaM, 0)
		}
		if win.Pressed(pixelgl.KeyRight) {
			p.direction = DirR
			move = pixel.V(deltaM, 0)
		}
		if win.Pressed(pixelgl.KeyDown) {
			p.direction = DirD
			move = pixel.V(0, -deltaM)
		}
		if win.Pressed(pixelgl.KeyUp) {
			p.direction = DirU
			move = pixel.V(0, deltaM)
		}

		if win.JustPressed(pixelgl.KeyV) {
			animInterval += 0.5
			fmt.Printf("animInterval %f\n", animInterval)
		}
		if win.JustPressed(pixelgl.KeyC) {
			animInterval -= 0.5
			fmt.Printf("animInterval %f\n", animInterval)
		}

		if win.JustPressed(pixelgl.KeyF) {
			speed += 0.5
			fmt.Printf("speed %f\n", speed)
		}
		if win.JustPressed(pixelgl.KeyD) {
			speed -= 0.5
			fmt.Printf("speed %f\n", speed)
		}

		if move.Len() > 0.0 {
			distance := elapsed.Seconds() * speed
			p.position = p.position.Add(move.Scaled(distance))
			p.animStep += distance
			if p.animStep > animInterval {
				p.anim = 1 + p.anim%(spriteCols-1)
				p.animStep = 0
			}
		} else {
			p.anim = 0
			p.animStep = 0
		}

		win.Clear(colornames.Burlywood)
		p.Draw(win)
		win.Update()

		last = now
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
