package pixel

import (
	"image"
	"image/color"
)

type Option struct {
	Threshold float32
	IncludeAA bool
}

func Match(img1, img2 image.Image, rect image.Rectangle, output *image.RGBA, opt *Option) int {
	var threshold float32
	threshold = 0.1
	if opt != nil {
		threshold = opt.Threshold
	}
	maxDelta := 35215 * threshold * threshold
	diff := 0

	includeAA := false
	if opt != nil {
		includeAA = opt.IncludeAA
	}

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			pos := image.Point{X: x, Y: y}
			delta := colorDelta(img1, img2, pos, pos, false)
			if delta > maxDelta {

				if !includeAA && (antialiased(img1, img2, x, y, rect.Max.X, rect.Max.Y) || antialiased(img2, img1, x, y, rect.Max.X, rect.Max.Y)) {

					if output != nil {
						drawPixel(output, pos, color.RGBA{255, 255, 0, 255})
					}

				} else {
					if output != nil {
						drawPixel(output, pos, color.RGBA{255, 0, 0, 255})
					}
					diff++
				}
			} else if output != nil {
				val := blend2(grayPixel(img1, pos), 0.1)
				v := uint8(val)
				drawPixel(output, pos, color.RGBA{v, v, v, 255})
			}
		}
	}
	return diff
}

func colorDelta(img1, img2 image.Image, k, m image.Point, yOnly bool) float32 {
	color1 := img1.At(k.X, k.Y)
	color2 := img2.At(m.X, m.Y)

	red1, green1, blue1, alpha1 := color1.RGBA()
	red1 = red1 >> 8
	green1 = green1 >> 8
	blue1 = blue1 >> 8
	alpha1 = alpha1 >> 8

	red2, green2, blue2, alpha2 := color2.RGBA()
	red2 = red2 >> 8
	green2 = green2 >> 8
	blue2 = blue2 >> 8
	alpha2 = alpha2 >> 8

	a1 := alpha1 / 255
	a2 := alpha2 / 255

	r1 := blend(red1, a1)
	g1 := blend(green1, a1)
	b1 := blend(blue1, a1)

	r2 := blend(red2, a2)
	g2 := blend(green2, a2)
	b2 := blend(blue2, a2)

	y := rgb2y(r1, g1, b1) - rgb2y(r2, g2, b2)
	if yOnly {
		return y
	}

	i := rgb2i(r1, g1, b1) - rgb2i(r2, g2, b2)
	q := rgb2q(r1, g1, b1) - rgb2q(r2, g2, b2)
	return 0.5053*y*y + 0.299*i*i + 0.1957*q*q
}

func blend(c, a uint32) float32 {
	return float32(255 + (c-255)*a)
}

func blend2(c, a float32) float32 {
	return 255 + (c-255)*a
}

func rgb2y(r, g, b float32) float32 {
	return r*0.29889531 + g*0.58662247 + b*0.11448223
}

func rgb2i(r, g, b float32) float32 {
	return r*0.59597799 - g*0.27417610 - b*0.32180189
}

func rgb2q(r, g, b float32) float32 {
	return r*0.21147017 - g*0.52261711 + b*0.31114694
}

func antialiased(imgA, imgB image.Image, x1, y1, w, h int) bool {
	x0 := max(x1-1, 0)
	y0 := max(y1-1, 0)
	x2 := min(x1+1, w-1)
	y2 := min(y1+1, h-1)

	pos := image.Point{X: x1, Y: y1}

	var zeroes, positives, negatives int
	var minX, minY, maxX, maxY int
	var _min, _max float32

	for x := x0; x <= x2; x++ {
		for y := y0; y <= y2; y++ {
			if x == x1 && y == y1 {
				continue
			}

			delta := colorDelta(imgA, imgA, pos, image.Point{X: x, Y: y}, true)
			if delta == 0 {
				zeroes++
			} else if delta > 0 {
				positives++
			} else if delta < 0 {
				negatives++
			}

			if zeroes > 2 {
				return false
			}

			if imgB == nil {
				continue
			}

			if delta < _min {
				_min = delta
				minX = x
				minY = y
			}
			if delta > _max {
				_max = delta
				maxX = x
				maxY = y
			}
		}
	}

	if imgB == nil {
		return true
	}

	if negatives == 0 || positives == 0 {
		return false
	}

	return (!antialiased(imgA, nil, minX, minY, w, h) && !antialiased(imgB, nil, minX, minY, w, h)) ||
		(!antialiased(imgA, nil, maxX, maxY, w, h) && !antialiased(imgB, nil, maxX, maxY, w, h))
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func grayPixel(img image.Image, pos image.Point) float32 {
	c := img.At(pos.X, pos.Y)
	red, green, blue, alpha := c.RGBA()
	red = red >> 8
	green = green >> 8
	blue = blue >> 8
	alpha = alpha >> 8

	a := alpha / 255
	r := blend(red, a)
	g := blend(green, a)
	b := blend(blue, a)
	return rgb2y(r, g, b)
}

func drawPixel(img *image.RGBA, pos image.Point, color color.Color) {
	img.Set(pos.X, pos.Y, color)
}
