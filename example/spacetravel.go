//taken from http://tech.nitoyon.com/en/blog/2016/01/18/space-travel-animated-gif/
package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"math"
	"math/rand"
)

var w, h float64 = 500, 250
var palette color.Palette = color.Palette{}
var zCycle float64 = 8
var zMin, zMax float64 = 1, 15

type Point struct {
	X, Y float64
}

type Circle struct {
	X, Y, Z, R float64
}

// Draw stars in order to generate perfect loop GIF
func (c *Circle) Draw(gc *draw2dimg.GraphicContext, ratio float64) {
	z := c.Z - ratio*zCycle

	for z < zMax {
		if z >= zMin {
			x, y, r := c.X/z, c.Y/z, c.R/z
			gc.SetFillColor(color.White)
			gc.Fill()
			draw2dkit.Circle(gc, w/2+x, h/2+y, r)
			gc.Close()
		}
		z += zCycle
	}
}

func drawFrame(circles []Circle, ratio float64) *image.Paletted {
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	gc := draw2dimg.NewGraphicContext(img)

	// Draw background
	gc.SetFillColor(color.Gray{0x11})
	draw2dkit.Rectangle(gc, 0, 0, w, h)
	gc.Fill()
	gc.Close()

	// Draw stars
	for _, circle := range circles {
		circle.Draw(gc, ratio)
	}

	// Dithering
	pm := image.NewPaletted(img.Bounds(), palette)
	draw.FloydSteinberg.Draw(pm, img.Bounds(), img, image.ZP)
	return pm
}

func spaceTravel(writer io.Writer) error {
	// Create 4000 stars
	circles := []Circle{}
	for len(circles) < 4000 {
		x, y := rand.Float64()*8-4, rand.Float64()*8-4
		if math.Abs(x) < 0.5 && math.Abs(y) < 0.5 {
			continue
		}
		z := rand.Float64() * zCycle
		circles = append(circles, Circle{x * w, y * h, z, 5})
	}

	// Intiialize palette (#000000, #111111, ..., #ffffff)
	palette = color.Palette{}
	for i := 0; i < 16; i++ {
		palette = append(palette, color.Gray{uint8(i) * 0x11})
	}

	// Generate 30 frames
	var images []*image.Paletted
	var delays []int
	count := 30
	for i := 0; i < count; i++ {
		pm := drawFrame(circles, float64(i)/float64(count))
		images = append(images, pm)
		delays = append(delays, 4)
	}

	// Output gif
	return gif.EncodeAll(writer, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
