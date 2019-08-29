package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"net/http"
	"strconv"
	"strings"
	"log"
	"fmt"
)

func main() {
	http.HandleFunc("/", imgHandler)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	queries := strings.Split(rawQuery, "&")
	queryMap := make(map[string]string)
	for _, q := range queries {
		item := strings.Split(q, "=")
		if len(item) < 2 {
			continue
		}
		queryMap[item[0]] = item[1]
	}
	x, _ := strconv.ParseFloat(queryMap["x"], 64)
	y, _ := strconv.ParseFloat(queryMap["y"], 64)
	ratio, _ := strconv.ParseFloat(queryMap["ratio"], 64)
	encodeImg(x, y, ratio, w)
}

func encodeImg(x, y, ratio float64, w http.ResponseWriter) {
	xmin := math.Min(x, x*ratio)
	xmax := math.Max(x, x*ratio)
	ymin := math.Min(y, y*ratio)
	ymax := math.Max(y, y*ratio)
	fmt.Println(xmin, xmax, ymin, ymax)
	// xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height := 1024, 1024
	// mycolor := color.RGBA{0xf3, 0x43, 0x11, 0xff}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			z := complex(x, y)
			avgcolor := imgAvgColor(img, px, py)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z, avgcolor))
		}
	}
	png.Encode(w, img) // NOTE: ignoring error
}

// imgAvgColor calculate avg subpixel color to reduce pixelation
func imgAvgColor(image *image.NRGBA, x, y int) color.RGBA {
	var r, g, b, a uint32
	left := image.At(x-1, y)
	right := image.At(x+1, y)
	up := image.At(x, y-1)
	down := image.At(x, y+1)
	leftr, leftg, leftb, lefta := left.RGBA()
	rightr, rightg, rightb, righta := right.RGBA()
	upr, upg, upb, upa := up.RGBA()
	downr, downg, downb, downa := down.RGBA()
	r = (leftr + rightr + upr + downr) / 4
	g = (leftg + rightg + upg + downg) / 4
	b = (leftb + rightb + upb + downb) / 4
	a = (lefta + righta + upa + downa) / 4
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func mandelbrot(z complex128, c color.RGBA) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			// return color.Gray{255 - contrast*n}
			return contrastColor(c, contrast, n)
		}
	}
	return c
}

func contrastColor(c color.RGBA, contrast uint8, n uint8) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, 255 - contrast*n}
}
