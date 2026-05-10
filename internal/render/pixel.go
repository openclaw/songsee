package render

import (
	"image"
	"image/color"
)

func setRGBA(img *image.RGBA, x, y int, c color.RGBA) {
	i := img.PixOffset(x, y)
	img.Pix[i+0] = c.R
	img.Pix[i+1] = c.G
	img.Pix[i+2] = c.B
	img.Pix[i+3] = c.A
}
