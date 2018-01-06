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
	DirNone
)

type animation struct {
	picture pixel.Picture
	sprite  *pixel.Sprite

	anim     int
	step     float64
	offset   int
	interval float64
}

func (a *animation) Update(dt time.Duration, distance float64, direction Direction) {
	if distance > 0 {
		a.offset = int(direction)*spriteCols + a.anim
		a.step += distance
		if a.step > a.interval {
			a.anim = 1 + a.anim%(spriteCols-1)
			a.step = 0
		}
	} else {
		a.offset = int(direction) * spriteCols
		a.anim = 0
		a.step = 0
	}
}

func (p *animation) Sprite() *pixel.Sprite {
	spriteWidth := p.picture.Bounds().Max.X / spriteCols
	spriteHeight := p.picture.Bounds().Max.Y / spriteRows

	spriteX := float64(p.offset%spriteCols) * spriteWidth
	spriteY := float64(p.offset/spriteCols) * spriteHeight

	bounds := pixel.R(spriteX, spriteY, spriteX+spriteWidth, spriteY+spriteHeight)
	if p.sprite == nil {
		p.sprite = pixel.NewSprite(p.picture, bounds)
	} else {
		p.sprite.Set(p.picture, bounds)
	}
	return p.sprite
}

type player struct {
	anim      animation
	direction Direction
	position  pixel.Vec
	speed     float64
}

const spriteCols = 9
const spriteRows = 4

func (p *player) Update(dt time.Duration, direction Direction) {
	if direction != DirNone {
		p.direction = direction
	}

	var move pixel.Vec
	switch direction {
	case DirL:
		move = pixel.V(-1, 0)
	case DirR:
		move = pixel.V(1, 0)
	case DirD:
		move = pixel.V(0, -1)
	case DirU:
		move = pixel.V(0, 1)
	}
	distance := move.Scaled(dt.Seconds() * p.speed)
	p.position = p.position.Add(distance)

	p.anim.Update(dt, distance.Len(), p.direction)
}

func (p *player) Draw(target pixel.Target) {
	p.anim.Sprite().Draw(target, pixel.IM.Moved(p.position))
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
		anim: animation{
			picture:  pic,
			interval: 4.5,
		},
		direction: DirR,
		speed:     76.0,
		position:  win.Bounds().Center(),
	}

	last := time.Now()
	for !win.Closed() {
		now := time.Now()
		elapsed := now.Sub(last)

		direction := DirNone
		if win.Pressed(pixelgl.KeyLeft) {
			direction = DirL
		}
		if win.Pressed(pixelgl.KeyRight) {
			direction = DirR
		}
		if win.Pressed(pixelgl.KeyDown) {
			direction = DirD
		}
		if win.Pressed(pixelgl.KeyUp) {
			direction = DirU
		}

		if win.JustPressed(pixelgl.KeyV) {
			p.anim.interval += 0.5
			fmt.Printf("animInterval %f\n", p.anim.interval)
		}
		if win.JustPressed(pixelgl.KeyC) {
			p.anim.interval -= 0.5
			fmt.Printf("animInterval %f\n", p.anim.interval)
		}

		if win.JustPressed(pixelgl.KeyF) {
			p.speed += 0.5
			fmt.Printf("speed %f\n", p.speed)
		}
		if win.JustPressed(pixelgl.KeyD) {
			p.speed -= 0.5
			fmt.Printf("speed %f\n", p.speed)
		}

		p.Update(elapsed, direction)

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
